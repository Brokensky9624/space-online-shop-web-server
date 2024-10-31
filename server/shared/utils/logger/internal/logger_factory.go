package internal

import (
	"fmt"
	"os"
	"time"

	"space.online.shop.web.server/shared/utils/logger/color"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// core
func GetFileCore(level zap.AtomicLevel, fileLogger *lumberjack.Logger) zapcore.Core {
	return zapcore.NewCore(
		NewMyLogEncoder(GetFileEncoderCfg()),
		zapcore.AddSync(fileLogger),
		level,
	)
}

func GetDefaultColorCore(level zap.AtomicLevel) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(GetDefaultColorCfg()),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

func GetColorCore(level zap.AtomicLevel, colorMap map[string]int) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(GetColorCfg(colorMap)),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

func GetStdCore(level zap.AtomicLevel) zapcore.Core {
	return zapcore.NewCore(
		NewMyLogEncoder(GetStdEncoderCfg()),
		zapcore.AddSync(os.Stdout),
		level,
	)
}

// encoder config
func GetFileEncoderCfg() zapcore.EncoderConfig {
	config := getDefaultEncoderCfg()
	config.EncodeTime = encodeFmtTime
	config.EncodeLevel = encodeFmtLevel
	config.EncodeCaller = encodeFmtCaller
	return config
}

func GetDefaultColorCfg() zapcore.EncoderConfig {
	config := getDefaultEncoderCfg()
	config.EncodeTime = encodeFmtTime
	config.EncodeLevel = encodeDefaultColorLevel
	config.EncodeCaller = encodeGoRoutineCaller
	return config
}

func GetColorCfg(colorMap map[string]int) zapcore.EncoderConfig {
	config := getDefaultEncoderCfg()
	config.EncodeTime = encodeFmtTime
	config.EncodeLevel = encodeColorLevel(colorMap)
	config.EncodeCaller = encodeGoRoutineCaller
	return config
}

func GetStdEncoderCfg() zapcore.EncoderConfig {
	config := getDefaultEncoderCfg()
	config.EncodeTime = encodeFmtTime
	config.EncodeLevel = encodeFmtLevel
	config.EncodeCaller = encodeFmtCaller
	return config
}

func getDefaultEncoderCfg() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
	}
}

// encode func
func encodeFmtTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encoder := zapcore.TimeEncoderOfLayout("2006/01/02-15:04:05.000")
	encoder(t, enc)
}

func encodeFmtLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	text := fmt.Sprintf("[%s]", l.CapitalString())
	enc.AppendString(text)
}

func encodeFmtCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	text := fmt.Sprintf("[%s]", caller.TrimmedPath())
	enc.AppendString(text)
}

func encodeDefaultColorLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	ct := getColorTextByLevel(l)
	ct.ApplyDefaultColorByLevel(l)
	enc.AppendString(ct.Build())
}

func encodeColorLevel(colorMap map[string]int) func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		ct := getColorTextByLevel(l)
		for _, isFg := range []bool{true, false} {
			var pos string
			if isFg {
				pos = "fg"
			} else {
				pos = "bg"
			}
			colorKey := fmt.Sprintf("%s_%s", getMapLevelKey(l), pos)
			code, ok := colorMap[colorKey]
			if !ok {
				code = 0x000000FF
			}
			colors := color.ParseLoggerColorCode(code)
			if colors == nil { // default color
				ct.ApplyDefaultColorByLevelAndPos(l, isFg)
			} else if len(colors) == 1 { // isSGR
				ct.SetSGRByPos(isFg, colors[0])
			} else if len(colors) == 3 { // isRGB
				ct.SetRGBByPos(isFg, [3]uint8{colors[0], colors[1], colors[2]})
			}
		}
		enc.AppendString(ct.Build())
	}
}

func encodeGoRoutineCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	goroutineText := fmt.Sprintf("[GO-%d]", getGoID())
	enc.AppendString(goroutineText)
	text := fmt.Sprintf("[%s]", caller.TrimmedPath())
	enc.AppendString(text)
}

// color
func getMapLevelKey(level zapcore.Level) string {
	var levelKey string
	switch level {
	case zapcore.DebugLevel:
		levelKey = "debug"
	case zapcore.InfoLevel:
		levelKey = "info"
	case zapcore.WarnLevel:
		levelKey = "warn"
	case zapcore.ErrorLevel:
		levelKey = "error"
	case zapcore.DPanicLevel:
		levelKey = "dpanic"
	case zapcore.PanicLevel:
		levelKey = "panic"
	case zapcore.FatalLevel:
		levelKey = "fatal"
	default:
		levelKey = "fatal"
	}
	return levelKey
}
func getColorTextByLevel(level zapcore.Level) *color.ColorText {
	return color.New(fmt.Sprintf("[%-5s]", level.CapitalString()))
}
