package types

import (
	"fmt"
	"reflect"
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

func ToProduct(input interface{}) Product {
	p := Product{}
	iVal := reflect.ValueOf(input)
	if iVal.Kind() == reflect.Struct {
		tool.CopyFields(&p, iVal)
	}
	return p
}

type CreateParam struct {
	Name         string `json:"name" required:"true"`
	Title        string `json:"title" required:"true"`
	Desc         string `json:"desc" required:"true"`
	Category     string `json:"category" required:"true"`
	Brand        string `json:"brand" required:"true"`
	Manufacturer string `json:"manufacturer" required:"true"`
	Status       uint   `json:"status"`
	Like         uint   `json:"like"`
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
		Like:         param.Like,
	}
}

func (param CreateParam) Check() error {
	if err := tool.CheckRequiredFields(param); err != nil {
		return err
	}
	if !slices.Contains(SupportedCategories, param.Category) {
		return fmt.Errorf("category %s not support", param.Category)
	}
	return nil
}

type DetailParam struct {
	ID uint `json:"id" required:"true"`
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
	ID uint `json:"id" required:"true"`
}

func (param LikeParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type DeleteParam struct {
	ID uint `json:"id" required:"true"`
}

func (param DeleteParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type DeleteBatchesParam struct {
	IDList []uint `json:"idList" required:"true"`
}

func (param DeleteBatchesParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type QueryParam struct {
	Title        string    `json:"title"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Brand        string    `json:"brand"`
	Manufacturer string    `json:"manufacturer"`
	Status       uint      `json:"status"`
	Like         uint      `json:"like"`
	UpdatedAt    time.Time `json:"updateAt"`
}

func (param QueryParam) Check() error {
	return tool.CheckRequiredFields(param)
}
