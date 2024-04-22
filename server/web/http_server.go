package web

import (
	"fmt"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/member"
	"space.online.shop.web.server/service"
)

var webServer *WebServer

type WebServer struct {
	Engine     *gin.Engine
	JWTMid     *jwt.GinJWTMiddleware
	SrvManager *service.ServiceManager
}

func New() *WebServer {
	mid, err := NewJWTMid()
	if err != nil {
		panic(fmt.Errorf("jwt err: %s", err))
	}
	webServer = &WebServer{
		Engine: gin.Default(),
		JWTMid: mid,
	}
	return webServer
}

func Server() *WebServer {
	return webServer
}

func (s *WebServer) SetSrvManager(srvManager *service.ServiceManager) *WebServer {
	s.SrvManager = srvManager
	return s
}

func (s *WebServer) Initialize() {
	s.RegisterRoute()
	go func() {
		webServer.Engine.Run(":3000")
	}()
}

func (s *WebServer) RegisterRoute() *WebServer {
	s.Engine.POST("/login", s.JWTMid.LoginHandler)
	s.Engine.POST("/refresh-token", s.JWTMid.RefreshHandler)
	apiGroup := s.Engine.Group("/api")
	apiGroup.Use(s.JWTMid.MiddlewareFunc())
	rootGroup := s.Engine.Group("/")
	member.NewREST(s.SrvManager, rootGroup).RegisterRoute()
	return s
}
