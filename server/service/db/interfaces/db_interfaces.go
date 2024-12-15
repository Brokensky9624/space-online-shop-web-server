package interfaces

import "gorm.io/gorm"

type OrderType string

const (
	ASC  OrderType = "asc"
	DESC OrderType = "desc"
)

type IDbService interface {
	GetDB() *gorm.DB
	Close() error
}

type IDbBuilder interface {
	BuildDB() (*gorm.DB, error)
}
