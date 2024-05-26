package model

import (
	"reflect"

	"gorm.io/gorm"
	"space.online.shop.web.server/util/tool"
)

type Product struct {
	gorm.Model
	Name         string `gorm:"size:50;not null"`
	Title        string `gorm:"size:256;not null"`
	Desc         string `gorm:"size:1024;not null"`
	Category     string `gorm:"size:50;not null"`
	Brand        string `gorm:"size:100;not null"`
	Manufacturer string `gorm:"size:200;not null"`
	Status       uint   `gorm:"size:128;not null"`
	Like         uint   `gorm:"size:1024000"`
	OwnerID      uint   `gorm:"not null"`
}

func (Product) TableName() string {
	return "product"
}

func (p *Product) SetID(productID uint) *Product {
	p.ID = productID
	return p
}

func (p *Product) SetOwner(userID uint) *Product {
	p.OwnerID = userID
	return p
}

func (p Product) IsOwner(userID uint) bool {
	return p.OwnerID == userID
}

func ToProductModel(input interface{}) Product {
	model := Product{}
	iValue := reflect.ValueOf(input)
	if iValue.Kind() == reflect.Struct {
		tool.CopyFields(&model, iValue)
	}
	return model
}
