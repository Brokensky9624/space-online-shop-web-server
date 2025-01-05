package web

import (
	"sync"

	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/member"
	"space.online.shop.web.server/rest/product"
	"space.online.shop.web.server/service"
	"space.online.shop.web.server/shared/utils/path"
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

func New(SrvManager *service.ServiceManager) *webServer {
	once.Do(func() {
		webSrv = &webServer{
			engine:     gin.Default(),
			SrvManager: SrvManager,
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

func (w *webServer) loadStaticFiles() {
	w.engine.StaticFile("/app.js", path.JoinRootPath("html", "app.js")) // 映射 app.js
	w.engine.Static("/css", path.JoinRootPath("html", "css"))           // 映射 css 文件夾
	w.engine.Static("/img", path.JoinRootPath("html", "img"))           // 映射 img 文件夾
	w.engine.Static("/js", path.JoinRootPath("html", "js"))             // 映射 js 文件夾（如果需要）

	w.engine.NoRoute(func(c *gin.Context) {
		c.File(path.JoinRootPath("html", "index.html"))
	})
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
	w.engine.GET("/auth", w.jwtAuth.GetAuthHandler())

	apiGroup := w.engine.Group("/api")
	apiGroup.Use(w.jwtAuth.GetMiddleware())
	memberREST := member.NewREST(w.SrvManager, apiGroup).RegisterRoute()
	w.engine.POST("/register", memberREST.Register)
	_ = product.NewREST(w.SrvManager, apiGroup).RegisterRoute()

	w.loadStaticFiles()
	return w
}
