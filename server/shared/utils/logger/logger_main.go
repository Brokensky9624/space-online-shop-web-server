package logger

import (
	"time"

	"space.online.shop.web.server/shared/utils/logger/internal"
)

var (
	STD    *internal.MyLogger
	SERVER *internal.MyLogger
)

func init() {
	STD = internal.NewMyLogger("std.properties", internal.WithPropertyUpdateInterval(3*time.Second))
	SERVER = internal.NewMyLogger("server.properties", internal.WithPropertyUpdateInterval(3*time.Second))
}

func GetLoggers() []*internal.MyLogger {
	return []*internal.MyLogger{STD, SERVER}
}

func RotateFileLoggers() {
	for _, l := range GetLoggers() {
		if l != nil {
			l.RotateFileLogger()
		}
	}
}

func CloseLoggers() {
	for _, l := range GetLoggers() {
		if l != nil {
			l.Close()
		}
	}
}
