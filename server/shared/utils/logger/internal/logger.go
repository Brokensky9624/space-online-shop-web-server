package internal

import (
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"space.online.shop.web.server/shared/interfaces"
	utilsPath "space.online.shop.web.server/shared/utils/path"

	"github.com/magiconair/properties"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const defaultPropertyUpdateInterval = 5 * time.Second

type MyLogger struct {
	wg                     sync.WaitGroup
	logger                 *zap.SugaredLogger
	fileLogger             *lumberjack.Logger
	atomicLevel            zap.AtomicLevel
	propertyPath           string
	property               MyLoggerProperty
	propertyUpdateInterval time.Duration
	closeCh                chan struct{}
	closeOnce              sync.Once
}

func (l *MyLogger) Debug(template string, args ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Debugf(template, args...)
}

func (l *MyLogger) Info(template string, args ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Infof(template, args...)
}

func (l *MyLogger) Warn(template string, args ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Warnf(template, args...)
}

func (l *MyLogger) Error(template string, args ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Errorf(template, args...)
}

func (l *MyLogger) Fatal(template string, args ...interface{}) {
	if l.logger == nil {
		l.init()
	}
	l.logger.Fatalf(template, args...)
}

func (l *MyLogger) RotateFileLogger() {
	if l.fileLogger != nil {
		l.fileLogger.Rotate()
	}
}

func (l *MyLogger) Close() {
	l.closeOnce.Do(func() {
		close(l.closeCh)
		l.wg.Wait()
	})
}

func (l *MyLogger) init() {
	l.atomicLevel = zap.NewAtomicLevel()
	l.createLogger()
	l.createUpdateThread()
}

func (l *MyLogger) createLogger() {
	p := loadProperty(l.propertyPath)
	l.assignLevel(p.level)

	var cores []zapcore.Core
	if p.isLogToFile {
		fileLogger := &lumberjack.Logger{
			Filename:   p.fileName,
			MaxSize:    p.fileMaxSize,
			MaxBackups: p.fileMaxBackups,
			MaxAge:     p.fileMaxAge,
		}
		cores = append(cores, GetFileCore(l.atomicLevel, fileLogger))
		l.fileLogger = fileLogger
	}
	if p.isLogToStd {
		if p.enableStdColor && p.isLogColorCustom {
			cores = append(cores, GetColorCore(l.atomicLevel, p.logColor))
		} else if p.enableStdColor && !p.isLogColorCustom {
			cores = append(cores, GetDefaultColorCore(l.atomicLevel))
		} else {
			cores = append(cores, GetStdCore(l.atomicLevel))
		}
	}
	if len(cores) == 0 { // if no core, at least non color std logger is needed
		cores = append(cores, GetStdCore(l.atomicLevel))
	}

	core := zapcore.NewTee(cores...)
	l.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func (l *MyLogger) assignLevel(levelText string) {
	zl := parseLevel(levelText)
	l.atomicLevel.SetLevel(zl)
}

func (l *MyLogger) createUpdateThread() {
	ticker := time.NewTicker(l.propertyUpdateInterval)
	l.wg.Add(1)
	go func() {
		defer func() {
			ticker.Stop()
			l.clear()
			l.wg.Done()
		}()
		for {
			select {
			case <-l.closeCh:
				return
			case <-ticker.C:
				l.update()
			}
		}
	}()
}

func (l *MyLogger) update() {
	var newConfig = loadProperty(l.propertyPath)
	ret := l.property.checkDifference(&newConfig)

	if ret == onlyLevelDiff {
		newLevel := parseLevel(newConfig.level)
		l.atomicLevel.SetLevel(newLevel)
	} else if ret == allDiff {
		l.clear()
		l.reload()
	}
}

func (l *MyLogger) reload() {
	l.createLogger()
}
func (l *MyLogger) clear() {
	l.logger.Sync()
	if l.fileLogger != nil {
		l.fileLogger.Close()
	}
	l.logger = nil
}

func (l *MyLogger) withOptions(opts ...interfaces.IOption[*MyLogger]) {
	for _, o := range opts {
		o.Apply(l)
	}
}

// = MyLoggerProperty
const (
	sameProperty int = iota
	onlyLevelDiff
	allDiff
)

type MyLoggerProperty struct {
	level            string
	isLogToStd       bool
	enableStdColor   bool
	isLogToFile      bool
	fileName         string
	fileMaxSize      int
	fileMaxBackups   int
	fileMaxAge       int
	isLogColorCustom bool
	logColor         map[string]int
}

func (p1 *MyLoggerProperty) checkDifference(p2 *MyLoggerProperty) (result int) {
	switch {
	case (p1.isLogToFile != p2.isLogToFile),
		(p1.isLogToStd != p2.isLogToStd),
		(p1.fileName != p2.fileName),
		(p1.fileMaxSize != p2.fileMaxSize),
		(p1.fileMaxBackups != p2.fileMaxBackups),
		(p1.fileMaxAge != p2.fileMaxAge):
		return allDiff
	case (p1.level != p2.level):
		return onlyLevelDiff
	}
	return sameProperty
}

// = export funcions
func NewMyLogger(propertyPath string, opts ...interfaces.IOption[*MyLogger]) *MyLogger {
	l := &MyLogger{
		propertyPath:           propertyPath,
		propertyUpdateInterval: defaultPropertyUpdateInterval,
		closeCh:                make(chan struct{}),
	}
	l.withOptions(opts...)
	return l
}

func WithPropertyUpdateInterval(interval time.Duration) interfaces.IOption[*MyLogger] {
	return interfaces.OptionFunc[*MyLogger](func(l *MyLogger) {
		l.propertyUpdateInterval = interval
	})
}

// private functions
func loadProperty(propertyPath string) MyLoggerProperty {
	var config MyLoggerProperty
	path := filepath.Join(utilsPath.LoggerDir, propertyPath)
	p, err := properties.LoadFile(path, properties.UTF8)
	if err != nil {
		return config
	}

	config.level = p.GetString("log.level", "FATAL")
	config.isLogToStd = p.GetBool("log.outToStd", false)
	config.enableStdColor = p.GetBool("log.outStd.color.enable", false)
	config.isLogToFile = p.GetBool("log.outToFile", false)
	config.fileName = filepath.Join(utilsPath.RootPath, p.GetString("log.file.filepath", "log/default_log/default.log"))
	config.fileMaxSize = p.GetInt("log.file.maxFileSize", 50)
	config.fileMaxBackups = p.GetInt("log.file.maxBackupNumber", 10)
	config.fileMaxAge = p.GetInt("log.file.maxBackupAge", 30)
	config.isLogColorCustom = p.GetBool("log.outStd.color.custom", false)
	config.logColor = map[string]int{}
	for _, l := range []string{"debug", "info", "warning", "error", "dpanic", "panic", "fatal"} {
		for _, f := range []string{"fg", "bg"} {
			s := p.GetString("log.outStd.color.custom."+l+"."+f, "0x000000FF")
			n, _ := strconv.ParseInt(s, 16, 32)
			config.logColor[l+"_"+f] = int(n)
		}
	}

	return config
}

func parseLevel(levelText string) zapcore.Level {
	switch levelText {
	case "DEBUG":
		return zap.DebugLevel
	case "INFO":
		return zap.InfoLevel
	case "WARN":
		return zap.WarnLevel
	case "ERROR":
		return zap.ErrorLevel
	case "FATAL":
		return zap.FatalLevel
	default:
		return zap.FatalLevel
	}
}
