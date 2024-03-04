package main

import (
	"space.online.shop.web.server/server/web"
)

func main() {
	web.Server().Initialize()
	select {}
}
