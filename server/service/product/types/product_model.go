package types

import (
	"fmt"
	"slices"
	"time"

	mysqlModel "space.online.shop.web.server/service/db/model"
	"space.online.shop.web.server/shared/utils/tool"
)

type ProductStatus uint

const (
	New ProductStatus = iota
	Preparing
	InStock
	OutofStock
)

var ( // FIXME: load config
	SupportedCategories = []string{
		"lamp",
		"chair",
	}
)

type Product struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
	Category     string    `json:"category"`
	Brand        string    `json:"brand"`
	Manufacturer string    `json:"manufacturer"`
	Status       uint      `json:"status"`
	Like         uint      `json:"like"`
	UpdatedAt    time.Time `json:"updateAt"`
	CreatedAt    time.Time `json:"createAt"`
}

type CreateParam struct {
	Name         string `json:"name" required:"true"`
	Title        string `json:"title" required:"true"`
	Desc         string `json:"desc" required:"true"`
	Category     string `json:"category" required:"true"`
	Brand        string `json:"brand" required:"true"`
	Manufacturer string `json:"manufacturer" required:"true"`
	Status       uint   `json:"status"`
}

func (param CreateParam) ToModel() mysqlModel.Product {
	return mysqlModel.Product{
		Name:         param.Name,
		Title:        param.Title,
		Desc:         param.Desc,
		Category:     param.Category,
		Brand:        param.Brand,
		Manufacturer: param.Manufacturer,
		Status:       param.Status,
	}
}

func (param CreateParam) Check() error {
	if err := tool.CheckRequiredFields(param); err != nil {
		return err
	}
	if !slices.Contains(SupportedCategories, param.Category) {
		return fmt.Errorf("category: `%s` not support", param.Category)
	}
	return nil
}

type DetailParam struct {
	ProductID uint `json:"productId" required:"true"`
}

func (param DetailParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type EditParam struct {
	ID           uint   `json:"id" required:"true"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Category     string `json:"category"`
	Brand        string `json:"brand"`
	Manufacturer string `json:"manufacturer"`
}

func (param EditParam) Check() error {
	return tool.CheckRequiredFields(param)
}

func (param EditParam) ToModel() mysqlModel.Product {
	return mysqlModel.Product{
		Name:         param.Name,
		Title:        param.Title,
		Desc:         param.Desc,
		Category:     param.Category,
		Brand:        param.Brand,
		Manufacturer: param.Manufacturer,
	}
}

type LikeParam struct {
	ProductID uint `json:"productId" required:"true"`
}

func (param LikeParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type DeleteParam struct {
	ProductID uint `json:"productId" required:"true"`
}

func (param DeleteParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type QueryParam struct {
	Title    string `json:"title"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Brand    string `json:"brand"`
	Page     int    `json:"page" required:"true"`
	PageSize int    `json:"pageSize" required:"true"`
}

func (param QueryParam) Offset() int {
	if param.Page < 1 || param.PageSize < 1 {
		return 0
	}
	return (param.Page - 1) * param.PageSize
}

func (param QueryParam) Limit() int {
	if param.Page < 1 || param.PageSize < 1 {
		return 0
	}
	return param.PageSize
}

func (param QueryParam) Check() error {
	return tool.CheckRequiredFields(param)
}

// type OrderParam struct {
// 	NameAsc      bool
// 	UpdatedAtAsc bool
// }

// type orderParamOption interface {
// 	apply(*OrderParam)
// }

// type orderParamOptionFunc func(*OrderParam)

// func (fn orderParamOptionFunc) apply(param *OrderParam) {
// 	fn(param)
// }

// func newDefaultOrderParam() *OrderParam {
// 	return &OrderParam{
// 		NameAsc: true,
// 	}
// }

// func NewOrderParam(opts ...orderParamOption) *OrderParam {
// 	param := newDefaultOrderParam()

// 	for _, opt := range opts {
// 		opt.apply(param)
// 	}

// 	return param
// }

// func WithOrderNameAsc(isAsc bool) orderParamOption {
// 	return orderParamOptionFunc(func(param *OrderParam) {
// 		param.NameAsc = isAsc
// 	})
// }

// func WithOrderUpdatedAtAsc(isAsc bool) orderParamOption {
// 	return orderParamOptionFunc(func(param *OrderParam) {
// 		param.UpdatedAtAsc = isAsc
// 	})
// }
