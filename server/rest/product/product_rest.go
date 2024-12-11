package product

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"space.online.shop.web.server/rest/response"
	"space.online.shop.web.server/service"
	memberTypes "space.online.shop.web.server/service/member/types"
	productTypes "space.online.shop.web.server/service/product/types"
)

type ProductREST struct {
	srvMngr        *service.ServiceManager
	apiRouterGroup *gin.RouterGroup
}

func NewREST(mngr *service.ServiceManager, routerGroup *gin.RouterGroup) *ProductREST {
	return &ProductREST{
		srvMngr:        mngr,
		apiRouterGroup: routerGroup,
	}
}

func (r *ProductREST) RegisterRoute() *ProductREST {
	productGroup := r.apiRouterGroup.Group("/product")
	{
		// single
		productGroup.POST("/create", r.Create)
		productGroup.PUT("/:id/edit", r.Edit)
		productGroup.PUT("/:id/like", r.Like)
		productGroup.GET("/:id", r.GetDetail)
		productGroup.DELETE("/:id/delete", r.Delete)
	}
	productsGroup := r.apiRouterGroup.Group("/products")
	{
		// batches
		productsGroup.GET("/query", r.Query)
	}
	return r
}

func (r *ProductREST) Create(c *gin.Context) {
	var params []productTypes.CreateParam
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

	srv := r.srvMngr.ProductSrv
	idList, err := srv.Create(member.ID, params...)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, response.SuccessRespObj(fmt.Sprintf("succeed in creating products, id list: %#v", idList), nil))
}

func (r *ProductREST) Edit(c *gin.Context) {
	idStr := c.Param("id")
	idUINT64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(idUINT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}

	var param productTypes.EditParam
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

	srv := r.srvMngr.ProductSrv
	if err := srv.Edit(member.ID, param); err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}

	msg := fmt.Sprintf("edit product %d succesfully!", id)
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, nil))
}

func (r *ProductREST) Like(c *gin.Context) {
	idStr := c.Param("id")
	idUINT64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(idUINT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}
	var param productTypes.LikeParam
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}
	param.ProductID = id

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
	srv := r.srvMngr.ProductSrv
	if err := srv.Like(member.ID, param); err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	msg := fmt.Sprintf("like product %d succesfully!", id)
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, nil))
}

func (r *ProductREST) Delete(c *gin.Context) {
	// idStr := c.Param("id")
	// idUINT64, err := strconv.ParseUint(idStr, 10, 64)
	// id := uint(idUINT64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, response.FailRespObj(err))
	// 	return
	// }
	// param := productTypes.DeleteParam{
	// 	ProductID: id,
	// }

	var params []productTypes.DeleteParam
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

	srv := r.srvMngr.ProductSrv
	idList, err := srv.Delete(member.ID, params...)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	c.JSON(http.StatusOK, response.SuccessRespObj(fmt.Sprintf("succeed in deleting products, id list: %#v", idList), nil))
}

func (r *ProductREST) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	idUINT64, err := strconv.ParseUint(idStr, 10, 64)
	id := uint(idUINT64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}
	param := productTypes.DetailParam{
		ProductID: id,
	}
	srv := r.srvMngr.ProductSrv
	product, err := srv.Detail(param)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}
	msg := fmt.Sprintf("get product %d detail succesfully!", id)
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, product))
}

func (r *ProductREST) Query(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}

	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailRespObj(err))
		return
	}

	search := c.DefaultQuery("search", "")

	qp := productTypes.QueryParam{
		Title:    search,
		Name:     search,
		Desc:     search,
		Brand:    search,
		Page:     int(page),
		PageSize: int(pageSize),
	}

	queryMap := c.Request.URL.Query()
	toOrderColumns := []string{}
	for k, v := range queryMap {
		if strings.HasSuffix(k, "_sort_order") {
			if slices.Contains(v, "asc") {
				toOrderColumns = append(toOrderColumns, strings.TrimSuffix(k, "_sort_order"))
			}
		}
	}

	srv := r.srvMngr.ProductSrv
	products, err := srv.Query(qp, toOrderColumns...)
	if err != nil {
		c.JSON(http.StatusOK, response.FailRespObj(err))
		return
	}

	msg := "query products succesfully!"
	c.JSON(http.StatusOK, response.SuccessRespObj(msg, products))
}
