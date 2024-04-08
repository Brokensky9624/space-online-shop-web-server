package rest

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/hello"
	"space.online.shop.web.server/rest/member"
)

func Register(r *gin.Engine) {
	authMiddware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "API",
		Key:         []byte("test key of space online"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*member.MemberParam); ok {
				return jwt.MapClaims{
					"id":       v.ID,
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals member.MemberParam
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			return "", nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"code": code, "message": message})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
	})
	if err != nil {
		panic("JWT Error " + err.Error())
	}

	r.POST("/login", authMiddware.LoginHandler)
	r.POST("/register")
	apiGroup := r.Group("/api")
	apiGroup.Use(authMiddware.MiddlewareFunc())
	{
		hello.Register(r)
		member.Register(r)
	}
}
