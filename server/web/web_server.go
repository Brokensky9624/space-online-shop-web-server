package web

import (
	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/server/rest"
)

var webServer *WebServer

type WebServer struct {
	Engine *gin.Engine
}

func init() {
	webServer = &WebServer{
		Engine: gin.Default(),
	}
	Register()
}

func Server() *WebServer {
	return webServer
}

func Register() {
	rest.Register(webServer.Engine)
}

func (s *WebServer) Initialize() {
	go func() {
		s.Engine.Run(":3000")
	}()
}
