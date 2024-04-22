package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"space.online.shop.web.server/service"
	"space.online.shop.web.server/service/db/mysql"
	"space.online.shop.web.server/service/member"
	"space.online.shop.web.server/web"
)

func main() {
	dbSrv, err := mysql.NewDBService()
	if err != nil {
		fmt.Println("mysql db build error:", err)
	}
	defer dbSrv.Close()

	// setup services to service manager
	memberSrv := member.NewService().SetDBService(dbSrv)
	srvManager := service.NewManager().SetMemberService(memberSrv)

	// setup web server and router
	web.New().SetSrvManager(srvManager).Initialize()

	// Set up signal handling to capture SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigCh:
		os.Exit(0)
	}
}
