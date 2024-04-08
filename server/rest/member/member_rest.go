package member

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/cmd/member"
	"space.online.shop.web.server/rest/response"
	"space.online.shop.web.server/util/tool"
)

func Register(r *gin.Engine) {
	r.POST("/register", func(c *gin.Context) {
		var user MemberParam
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// cmd register
		if err := member.Register(tool.StructToMap(user)); err != nil {
			c.JSON(http.StatusOK, response.FailRespObj(err))
			return
		}
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Register member " + user.Username + "!",
		})
	})
	r.POST("/login", func(c *gin.Context) {
		var user MemberParam
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// cmd create
		if err := member.Register(tool.StructToMap(user)); err != nil {
			c.JSON(http.StatusOK, response.FailRespObj(err))
			return
		}
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Create member " + user.Username + "!",
		})
	})

	group := r.Group("/api")
	group.PUT("/edit", func(c *gin.Context) {
		var user MemberEditParam
		// Get the "name" parameter from the route
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// cmd edit
		if err := member.Edit(tool.StructToMap(user)); err != nil {
			c.JSON(http.StatusOK, response.FailRespObj(err))
			return
		}
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Edit member " + user.Username + "!",
		})
	})
	group.DELETE("/delete/:username", func(c *gin.Context) {
		// Get the "name" parameter from the route
		username := c.Param("username")
		// cmd delete
		if err := member.Delete(username); err != nil {
			c.JSON(http.StatusOK, response.FailRespObj(err))
			return
		}
		// Set the response
		c.JSON(http.StatusOK, gin.H{
			"message": "Delete member " + username + "!",
		})
	})
	group.GET("/list", func(c *gin.Context) {
		// cmd list
		members, err := member.List()
		if err != nil {
			c.JSON(http.StatusOK, response.FailRespObj(err))
			return
		}
		ret := make([]interface{}, 0)
		for _, m := range members {
			ret = append(ret, m)
		}
		c.JSON(http.StatusOK, response.SuccessRespObj(ret))
	})
}
