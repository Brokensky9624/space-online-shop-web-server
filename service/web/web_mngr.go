package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json: "id"`
	Username string `json: "username"`
	Email    string `json: "email"`
}

var webMngr *WebMngr

type WebMngr struct {
	Engine *gin.Engine
}

func New() *WebMngr {
	webMngr = &WebMngr{}
	return webMngr
}

func (mngr *WebMngr) Init() {
	r := gin.Default()
	mngr.Engine = r
	mngr.Register()
	go func() {
		r.Run(":3000")
	}()
}

func (mngr *WebMngr) Register() {
	r := mngr.Engine

	api := r.Group("/api")
	{
		hello := api.Group("/hello")
		{
			hello.GET("/:name", func(c *gin.Context) {
				// Get the "name" parameter from the route
				name := c.Param("name")

				// Set the response
				c.JSON(http.StatusOK, gin.H{
					"message": "Hello " + name + "!",
				})
			})
		}

		member := api.Group("/member")
		{
			member.POST("register", func(c *gin.Context) {
				var user User
				// Get the "name" parameter from the route
				if err := c.ShouldBind(&user); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
				}
				// Set the response
				c.JSON(http.StatusOK, gin.H{
					"message": "Register member " + user.Username + "!",
				})
			})
		}
	}

	// r.GET("/api/hello/:name", func(c *gin.Context) {
	// 	// Get the "name" parameter from the route
	// 	name := c.Param("name")

	// 	// Set the response
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "Hello " + name + "!",
	// 	})
	// })
}
