package main

import (
	"os"
	"os/signal"
	"syscall"

	mysqlSrv "space.online.shop.web.server/service/db/mysql"
	webSrvr "space.online.shop.web.server/web"
)

func main() {
	webSrvr.New().Initialize()
	_, err := mysqlSrv.NewDB()
	if err != nil {
		// log fail
	}
	defer mysqlSrv.CloseDB()
	// Set up signal handling to capture SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		os.Exit(0)
	}
}
