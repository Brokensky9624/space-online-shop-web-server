package web

import (
	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest"
)

var webServer *WebServer

type WebServer struct {
	Engine *gin.Engine
}

func New() *WebServer {
	webServer = &WebServer{
		Engine: gin.Default(),
	}
	return webServer
}

func Server() *WebServer {
	return webServer
}

func Register() {
	rest.Register(webServer.Engine)
}

func (s *WebServer) Initialize() {
	Register()
	go func() {
		webServer.Engine.Run(":3000")
	}()
}
