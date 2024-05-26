package member

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/response"
	"space.online.shop.web.server/service"
	MemberTypes "space.online.shop.web.server/service/member/types"
)

func NewREST(srvManager *service.ServiceManager, routerGroup *gin.RouterGroup) *MemberREST {
	return &MemberREST{
		SrvManager: srvManager,
		apiGroup:   routerGroup,
	}
}

type MemberREST struct {
	apiGroup   *gin.RouterGroup
	SrvManager *service.ServiceManager
}

func (r *MemberREST) RegisterRoute() *MemberREST {
	memberGroup := r.apiGroup.Group("/member")
	{
		memberGroup.PUT("/edit", r.Edit)
		memberGroup.DELETE("/delete/:account", r.Delete)
		memberGroup.GET("/list", r.Members)
		memberGroup.GET("/:account", r.Member)
	}
	return r
}

func (r *MemberREST) Register(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	var user MemberTypes.MemberCreateParam
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// register member
	if err := memberSrv.Create(user); err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Register member %s successful !", user.Account)
	c.JSON(http.StatusOK, response.SuccessRespObj(message, nil))
}

func (r *MemberREST) Edit(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	var user MemberTypes.MemberEditParam
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// edit member
	if err := memberSrv.Edit(user); err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Edit member %s successful !", user.Account)
	c.JSON(http.StatusOK, response.SuccessRespObj(message, nil))
}

func (r *MemberREST) Delete(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	account := c.Param("account")
	// delete member
	if err := memberSrv.Delete(MemberTypes.MemberDeleteParam{
		Account: account,
	}); err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	message := fmt.Sprintf("Delete member %s successful !", account)
	c.JSON(http.StatusOK, response.SuccessRespObj(message, nil))
}

func (r *MemberREST) Members(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	// get member list
	memberList, err := memberSrv.Members()
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	memberListLen := len(memberList)
	dataList := make([]interface{}, memberListLen)
	for i, member := range memberList {
		dataList[i] = member
	}
	c.JSON(http.StatusOK, response.SuccessRespObj("", dataList...))
}

func (r *MemberREST) Member(c *gin.Context) {
	memberSrv := r.SrvManager.MemberService()
	account := c.Param("account")
	// query member
	member, err := memberSrv.Member(MemberTypes.MemberInfoParam{
		Account: account,
	})
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, response.SuccessRespObj("", member))
}
