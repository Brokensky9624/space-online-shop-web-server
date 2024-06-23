package interfaces

import "gorm.io/gorm"

type IDbService interface {
	GetDB() *gorm.DB
	Close() error
}

type IDbBuilder interface {
	BuildDB() (*gorm.DB, error)
}
