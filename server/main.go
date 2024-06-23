package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"space.online.shop.web.server/service/common"
	"space.online.shop.web.server/service/db"
	"space.online.shop.web.server/service/db/builder"
	"space.online.shop.web.server/util/logger"
	_ "space.online.shop.web.server/util/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var stoppers []common.Stoppable
	defer func() {
		cancel()
		common.StopAll(stoppers...)
	}()
	// prepare db
	dbBuilder := builder.NewMysqlDbBuilder(
		"space_online_admin",
		"space_online_is_666",
		"masterDB",
		builder.WithCharset("utf8mb4"),
		builder.WithParseTime("true"),
		builder.WithLoc("Local"),
		builder.WithTimeout("10s"),
		builder.WithReadTimeout("30s"),
		builder.WithWriteTimeout("30s"),
		builder.WithMaxIdleConns(20),
		builder.WithMaxOpenConns(200),
		builder.WithConnMaxIdleTime(20*time.Minute),
	)

	dbSrv := db.NewDbService(ctx, dbBuilder)
	stoppers = append(stoppers, dbSrv)
	go dbSrv.Run()

	// setup services to service manager
	// memberSrv := member.NewService(dbSrv)
	// productSrv := product.NewService(dbSrv)
	// srvManager := service.NewManager().
	// 	SetMemberService(memberSrv).
	// 	SetProductService(productSrv)

	// setup web server and router
	// web.New().SetSrvManager(srvManager).Initialize()

	logger.SERVER.Debug("all ready")
	name := "Jason"
	logger.SERVER.Debug("%s did good job.", name)

	// Set up signal handling to capture SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		c := <-sigCh
		switch c {
		case syscall.SIGHUP:
			// rotate logger
			logger.FileLoggerRotate()
		case syscall.SIGINT, syscall.SIGTERM:
			return
		}
	}
}
