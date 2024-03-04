package member

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type member struct {
	Username string `json: "username"`
	Email    string `json: "email"`
	Password string `json: "password"`
}

func Register(r *gin.Engine) {
	group := r.Group("/cmd/member")
	group.POST("/register", func(c *gin.Context) {
		var user member
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
	group.DELETE("/delete/:memberID", func(c *gin.Context) {
		// Get the "name" parameter from the route
		memberID := c.Param("memberID")

		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Delete member " + memberID + "!",
		})
	})
	group.PUT("/edit", func(c *gin.Context) {
		var user member
		// Get the "name" parameter from the route
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Edit member " + user.Username + "!",
		})
	})
	group.GET("/list", func(c *gin.Context) {
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "list all members!",
		})
	})
}
