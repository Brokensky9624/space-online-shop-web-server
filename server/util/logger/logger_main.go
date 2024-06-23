package logger

import "space.online.shop.web.server/util/logger/internal"

var STD = internal.NewMyLogger("cfg/logger/std.properties")
var SERVER = internal.NewMyLogger("cfg/logger/server.properties")

func FileLoggerRotate() {
	STD.FileLoggerRotate()
	SERVER.FileLoggerRotate()
}

func FileLoggerClose() {
	STD.FileLoggerClose()
	SERVER.FileLoggerClose()
}
