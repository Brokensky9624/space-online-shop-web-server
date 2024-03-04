package main

import (
	"fmt"
	"space-online-shop-web-server/service/web"
)

func main() {
	fmt.Println("666")
	mngr := web.New()
	mngr.Init()
	select {}
}
