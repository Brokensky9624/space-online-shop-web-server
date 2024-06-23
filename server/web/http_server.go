package web

import (
	"sync"

	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/member"
	"space.online.shop.web.server/rest/product"
	"space.online.shop.web.server/service"
	"space.online.shop.web.server/web/jwt"
)

var (
	once   sync.Once
	webSrv *webServer
)

type webServer struct {
	engine     *gin.Engine
	jwtAuth    jwt.IJWTAuth
	SrvManager *service.ServiceManager
}

func New() *webServer {
	once.Do(func() {
		webSrv = &webServer{
			engine: gin.Default(),
		}
	})
	webSrv.prepare()
	return webSrv
}

func Server() *webServer {
	return webSrv
}

func (w *webServer) prepare() {
	if w.jwtAuth == nil {
		w.loadJWTAuth()
	}
}

func (w *webServer) loadJWTAuth() {
	f := jwt.GetJWTFactory()
	w.jwtAuth = f.GetJWTAuth()
}

func (w *webServer) Initialize() {
	w.RegisterRoute()
	go func() {
		w.engine.Run(":3000")
	}()
}

func (w *webServer) RegisterRoute() *webServer {
	w.engine.POST("/login", w.jwtAuth.GetLoginHandler())
	w.engine.POST("/refresh-token", w.jwtAuth.GetRefreshHandler())
	apiGroup := w.engine.Group("/api")
	apiGroup.Use(w.jwtAuth.GetMiddleware())
	memberREST := member.NewREST(w.SrvManager, apiGroup).RegisterRoute()
	w.engine.POST("/register", memberREST.Register)
	_ = product.NewREST(w.SrvManager, apiGroup).RegisterRoute()
	return w
}
