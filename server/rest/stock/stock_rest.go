package stock

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/parse"
	"space.online.shop.web.server/rest/response"
	"space.online.shop.web.server/service"

	"space.online.shop.web.server/service/db/model"
	dbTypes "space.online.shop.web.server/service/db/types"
	memberTypes "space.online.shop.web.server/service/member/types"
	stockTypes "space.online.shop.web.server/service/stock/types"
)

type StockREST struct {
	srvMngr        *service.ServiceManager
	apiRouterGroup *gin.RouterGroup
}

func NewREST(mngr *service.ServiceManager, routerGroup *gin.RouterGroup) *StockREST {
	return &StockREST{
		srvMngr:        mngr,
		apiRouterGroup: routerGroup,
	}
}

func (r *StockREST) RegisterRoute() *StockREST {
	stockGroup := r.apiRouterGroup.Group("/stock")
	{
		stockGroup.POST("/create", r.Create)
		stockGroup.PUT("/:id/edit", r.Edit)
		stockGroup.GET("/:id", r.Detail)
		stockGroup.DELETE("/delete", r.Delete)
	}
	stocksGroup := r.apiRouterGroup.Group("/stocks")
	{
		stocksGroup.GET("/query", r.Query)
	}
	return r
}

func (r *StockREST) Create(c *gin.Context) {
	var params []stockTypes.CreateParam

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}

	user, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.FailRespObj(errors.New("Unauthorized")))
		return
	}
	member, ok := user.(*memberTypes.Member)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.FailRespObj(errors.New("internal Server Error")))
		return
	}

	srv := r.srvMngr.StockSrv
	idList, err := srv.Create(member.ID, params...)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, response.SuccessRespObj(fmt.Sprintf("succeed in creating stocks, id list: %#v", idList), nil))
}

func (r *StockREST) Edit(c *gin.Context) {
	idStr := c.Param("id")
	idUINT64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(idUINT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}

	var param stockTypes.EditParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}
	param.ID = id

	user, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.FailRespObj(errors.New("Unauthorized")))
		return
	}
	member, ok := user.(*memberTypes.Member)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.FailRespObj(errors.New("internal Server Error")))
		return
	}

	srv := r.srvMngr.StockSrv
	if err := srv.Edit(member.ID, param); err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}

	msg := fmt.Sprintf("edit stock %d successfully!", id)
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, nil))
}

func (r *StockREST) Delete(c *gin.Context) {
	var params []stockTypes.DeleteParam
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}

	// get user
	user, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.FailRespObj(errors.New("Unauthorized")))
		return
	}
	member, ok := user.(*memberTypes.Member)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.FailRespObj(errors.New("internal Server Error")))
		return
	}

	srv := r.srvMngr.StockSrv
	idList, err := srv.Delete(member.ID, params...)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, response.SuccessRespObj("succeed in deleting stocks", idList))
}

func (r *StockREST) Detail(c *gin.Context) {
	idStr := c.Param("id")
	idUINT64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(idUINT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}
	param := stockTypes.DetailParam{
		StockID: id,
	}
	srv := r.srvMngr.StockSrv
	stock, err := srv.Detail(param)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	msg := fmt.Sprintf("get stock %d detail successfully!", id)
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, stock))
}

func (r *StockREST) Query(c *gin.Context) {
	queryParser := parse.NewQueryParser(
		c.Request.URL.Query(),
		model.StockColumnMap(),
	)

	srv := r.srvMngr.StockSrv
	stocks, err := srv.Query(
		dbTypes.ConcatConditions(
			queryParser.Counter().Conditions(),
			queryParser.SortOrder().Conditions(),
			queryParser.Searcher().Conditions(),
		)...,
	)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}

	msg := "query stocks successfully!"
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, stocks))
}
