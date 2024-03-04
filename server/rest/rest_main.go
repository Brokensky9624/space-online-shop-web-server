package rest

import (
	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/server/rest/hello"
	"space.online.shop.web.server/server/rest/member"
)

func Register(r *gin.Engine) {
	hello.Register(r)
	member.Register(r)
}
