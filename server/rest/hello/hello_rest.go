package hello

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	group := r.Group("/cmd/hello")
	group.GET("/:name", func(c *gin.Context) {
		// Get the "name" parameter from the route
		name := c.Param("name")
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + name + "!",
		})
	})
}
