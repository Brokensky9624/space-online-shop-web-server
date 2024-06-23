package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/magiconair/properties"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"space.online.shop.web.server/cfg"
)

type MyLogger struct {
	propertyPath string
	atomicLevel  zap.AtomicLevel
	config       *loggerProperties
	fileLogger   *lumberjack.Logger
	logger       *zap.SugaredLogger
}

const (
	MONITOR_PERIOD = 5 * time.Second
)

const (
	sameProperties = iota
	onlyLevelDiff
	allDiff
)

type loggerProperties struct {
	level            string
	filePath         string
	maxFileSize      int
	maxAge           int
	maxBackupNumber  int
	localTime        bool
	isStdColorEnable bool
	isStdColorCustom bool
	outToStd         bool
	outToFile        bool
}

func (p *loggerProperties) diffProperties(np *loggerProperties) int {
	switch {
	case p.filePath != np.filePath,
		p.maxFileSize != np.maxFileSize,
		p.maxAge != np.maxAge,
		p.maxBackupNumber != np.maxBackupNumber,
		p.localTime != np.localTime,
		p.isStdColorEnable != np.isStdColorEnable,
		p.isStdColorCustom != np.isStdColorCustom,
		p.outToStd != np.outToStd,
		p.outToFile != np.outToFile:
		return allDiff
	case p.level != np.level:
		return onlyLevelDiff
	default:
		return sameProperties
	}
}

// constructer
func NewMyLogger(propertyPath string) *MyLogger {
	l := &MyLogger{
		propertyPath: propertyPath,
		atomicLevel:  zap.NewAtomicLevel(),
	}
	l.init()
	return l
}

// export functions

func (l *MyLogger) Info(msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

func (l *MyLogger) Debug(msg string, args ...any) {
	l.logger.Debug(fmt.Sprintf(msg, args...))
}

func (l *MyLogger) Warn(msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...))
}

func (l *MyLogger) Error(msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...))
}

func (l *MyLogger) Panic(msg string, args ...any) {
	l.logger.Panic(fmt.Sprintf(msg, args...))
}

func (l *MyLogger) Fatal(msg string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(msg, args...))
}

func (l *MyLogger) FileLoggerRotate() {
	if l.fileLogger != nil {
		l.fileLogger.Rotate()
	}
}

func (l *MyLogger) FileLoggerClose() {
	if l.fileLogger != nil {
		l.fileLogger.Close()
		l.fileLogger = nil
	}
}

// private
func (l *MyLogger) init() {
	l.initLogger()
	l.monitorThreadRun()

}

func (l *MyLogger) clear() {
	if l.logger != nil {
		l.logger.Sync()
	}
	l.FileLoggerClose()
}

func (l *MyLogger) reload() {
	l.initLogger()
}

func (l *MyLogger) initLogger() {
	// load properties and set level
	l.config = genLoggerProperties(l.propertyPath)
	l.atomicLevel.SetLevel(parseToZapCoreLevel(l.config.level))

	// prepare write syncer
	var writeSyncerList []zapcore.WriteSyncer
	if l.config.outToFile {
		l.fileLogger = genFileLogger(l.config)
		writeToFileSyncer := zapcore.AddSync(l.fileLogger)
		writeSyncerList = append(writeSyncerList, writeToFileSyncer)
	}
	if l.config.outToStd {
		if l.config.isStdColorEnable {

		} else {
			stdSyncer := zapcore.AddSync(os.Stdout)
			writeSyncerList = append(writeSyncerList, stdSyncer)
		}
	}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:    "M",
		LevelKey:      "L",
		TimeKey:       "T",
		NameKey:       "N",
		CallerKey:     "C",
		StacktraceKey: "S",
		EncodeTime:    encodeFmtTime,
		EncodeLevel:   encodeFmtLevel,
		EncodeCaller:  encodeFmtCaller,
	}
	core := zapcore.NewCore(
		NewSpaceOnlineEncoder(encoderCfg),
		zapcore.NewMultiWriteSyncer(writeSyncerList...),
		l.atomicLevel,
	)
	l.logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func (l *MyLogger) monitorThreadRun() {
	go func() {
		for {
			time.Sleep(MONITOR_PERIOD)
			l.detectPropertiesChange()
		}
	}()
}

func (l *MyLogger) detectPropertiesChange() {
	newProperties := genLoggerProperties(l.propertyPath)
	ret := l.config.diffProperties(newProperties)
	if ret == onlyLevelDiff {
		newLevel := parseToZapCoreLevel(newProperties.level)
		l.atomicLevel.SetLevel(newLevel)
	} else if ret == allDiff {
		l.clear()
		l.reload()
	}
}

func genLoggerProperties(propertyPath string) *loggerProperties {
	p, _ := properties.LoadFile(propertyPath, properties.UTF8)
	fp := cfg.JoinRootPath(p.GetString("log.file.filepath", "log/defautl_log/default_log.log"))
	return &loggerProperties{
		level:            p.GetString("log.level", "ERROR"),
		filePath:         fp,
		maxFileSize:      p.GetInt("log.file.maxFileSize", 50),
		maxBackupNumber:  p.GetInt("log.file.maxBackupNumber", 10),
		maxAge:           p.GetInt("log.file.maxAge", 10),
		outToFile:        p.GetBool("log.outToFile", false),
		outToStd:         p.GetBool("log.outToStd", false),
		isStdColorEnable: p.GetBool("log.outStd.color.enable", false),
		isStdColorCustom: p.GetBool("log.outStd.color.custom", false),
		localTime:        p.GetBool("log.file.localTime", false),
	}
}

func genFileLogger(config *loggerProperties) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.filePath,
		MaxSize:    config.maxFileSize,
		MaxAge:     config.maxAge,
		MaxBackups: config.maxBackupNumber,
		LocalTime:  config.localTime,
	}
}

func parseToZapCoreLevel(level string) zapcore.Level {
	switch level {
	case "INFO":
		return zap.InfoLevel
	case "DEBUG":
		return zap.DebugLevel
	case "WARN":
		return zap.WarnLevel
	case "ERROR":
		return zap.ErrorLevel
	case "FATAL":
		return zap.FatalLevel
	case "PANIC":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
