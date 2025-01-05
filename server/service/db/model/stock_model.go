package model

import (
	"sync"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/db"
)

var (
	stockOnce      sync.Once
	stockColumnMap map[string]struct{}
)

func FetchStockColumnMap() {
	stockOnce.Do(func() {
		stockColumnMap = make(map[string]struct{})
		db := db.Service()
		stmt := gorm.Statement{DB: db.DB}
		if err := stmt.Parse(&Stock{}); err != nil {
			panic(err)
		}

		for _, field := range stmt.Schema.Fields {
			if field.DBName != "" {
				stockColumnMap[field.DBName] = struct{}{}
			}
		}
	})
}

func StockColumnMap() map[string]struct{} {
	if stockColumnMap == nil {
		FetchStockColumnMap()
	}
	return stockColumnMap
}

type Stock struct {
	gorm.Model
	ProductID   uint
	Quantity    uint    `gorm:"size:128;not null"`
	Price       float64 `gorm:"size:128;not null"`
	Description string  `gorm:"size:200;not null"`
}

func (Stock) TableName() string {
	return "stock"
}
