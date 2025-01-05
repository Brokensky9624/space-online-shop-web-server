package types

import (
	"time"

	mysqlModel "space.online.shop.web.server/service/db/model"
)

type Stock struct {
	ID          uint
	ProductID   uint
	Quantity    uint      `json:"quantity"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateParam struct {
	ProductID   uint    `json:"productId" required:"true"`
	Quantity    uint    `json:"quantity"`
	Price       float64 `json:"price" required:"true"`
	Description string  `json:"description" required:"true"`
}

func (param CreateParam) ToModel() mysqlModel.Stock {
	return mysqlModel.Stock{
		ProductID:   param.ProductID,
		Quantity:    param.Quantity,
		Price:       param.Price,
		Description: param.Description,
	}
}

type DetailParam struct {
	StockID uint `json:"stockId" required:"true"`
}

type EditParam struct {
	ID          uint    `json:"id" required:"true"`
	Quantity    uint    `json:"quantity"`
	Price       float64 `json:"price" required:"true"`
	Description string  `json:"description" required:"true"`
}

func (param EditParam) ToModel() mysqlModel.Stock {
	return mysqlModel.Stock{
		Quantity:    param.Quantity,
		Price:       param.Price,
		Description: param.Description,
	}
}

type DeleteParam struct {
	StockID uint `json:"stockId" required:"true"`
}

type QueryParam struct {
	ProductID uint `json:"productId"`
	Page      int  `json:"page" required:"true"`
	PageSize  int  `json:"pageSize" required:"true"`
}
