package web

import (
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/service"
	memberTypes "space.online.shop.web.server/service/member/types"
)

const (
	identityKey  = "id"                               // indetiyKey for JWT claim
	identityRole = "role"                             // identityRole for JWT claim
	secretKey    = "5BYrir4vrBMB2oFJVywHFSrvlim6kCRn" // secret key for JWT encrypt
)

func NewJWTMid() (*jwt.GinJWTMiddleware, error) {
	memberSrv := service.Manager().MemberService()
	// JWT middleware initialization
	mid, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "test zone",        // it used to indicate proteced area or resource in server.
		SigningAlgorithm: "HS256",            // specific cryptographic algorithm used to sign the JWT tokens, default is `HS256`
		Key:              []byte(secretKey),  // secret key to generate or verify JWT tokens
		Timeout:          time.Hour,          //  expiration time of JWT tokens, after this time, need relogin or refresh token
		MaxRefresh:       time.Hour * 24 * 7, // expiration time of refresh JWT token by old JWT token, after this time, need relogin
		IdentityKey:      identityKey,        // this key used to store User identify information in JWT token claims
		PayloadFunc: func(data interface{}) jwt.MapClaims { // this function used to generate JWT token claims by User data
			if v, ok := data.(memberTypes.Member); ok {
				return jwt.MapClaims{
					identityKey:  v.ID,
					identityRole: v.Role,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} { // this function used to extract JWT token claims from restful request and generate User data
			claims := jwt.ExtractClaims(c)
			role := 0
			if v, ok := claims[identityRole].(float64); ok {
				role = int(v)
			}
			return memberTypes.Member{
				ID:   claims[identityKey].(string),
				Role: memberTypes.MemberRole(role),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) { // this function used to authenticate login user
			var loginUser memberTypes.Member
			if err := c.ShouldBindJSON(&loginUser); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			matchUser, err := memberSrv.AuthAndMember() // check user in user list, it might query database in real scenario
			if err != nil {
				return nil, err
			}
			if matchUser != nil {
				return &memberTypes.Member{
					ID:       matchUser.ID,
					Username: matchUser.Username,
					Password: matchUser.Password,
					Role:     matchUser.Role,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool { // this function used to check user authorization
			if u, ok := data.(*memberTypes.Member); ok {
				if strings.HasPrefix(c.Request.URL.Path, "/api/protect") {
					return u.Role == memberTypes.Admin
				}
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	return mid, err
}
