package model

import (
	"sync"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/db"
)

var (
	once             sync.Once
	productColumnMap map[string]struct{}
)

func FetchProductColumnMap() {
	once.Do(func() {
		productColumnMap = make(map[string]struct{})
		db := db.Service()
		stmt := gorm.Statement{DB: db.DB}
		if err := stmt.Parse(&Product{}); err != nil {
			panic(err)
		}

		for _, field := range stmt.Schema.Fields {
			if field.DBName != "" {
				productColumnMap[field.DBName] = struct{}{}
			}
		}
	})
}

func ProductColumnMap() map[string]struct{} {
	if productColumnMap == nil {
		FetchProductColumnMap()
	}
	return productColumnMap
}

type Product struct {
	gorm.Model
	Name         string   `gorm:"size:50;not null"`
	Title        string   `gorm:"size:256;not null"`
	Description  string   `gorm:"size:1024;not null"`
	Category     string   `gorm:"size:50;not null"`
	Brand        string   `gorm:"size:100;not null"`
	Manufacturer string   `gorm:"size:200;not null"`
	Status       uint     `gorm:"size:128;not null"`
	LikedBy      []Member `gorm:"many2many:member_product_likes;" json:"-"`
}

func (Product) TableName() string {
	return "product"
}

func (p *Product) SetID(productID uint) *Product {
	p.ID = productID
	return p
}
