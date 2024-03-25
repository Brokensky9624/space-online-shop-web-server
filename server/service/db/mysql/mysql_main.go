package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"space.online.shop.web.server/service/db/model"
)

type dbCfg struct {
	ip           string
	port         string
	userName     string
	password     string
	databaseName string
}

var ormDB *gorm.DB

func NewDB() (*gorm.DB, error) {
	cfg := &dbCfg{
		ip:           "localhost",
		port:         "3306",
		userName:     "space_online_admin",
		password:     "space_online_is_666",
		databaseName: "masterDB",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.userName, cfg.password, cfg.ip, cfg.port, cfg.databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}
	// ping db
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	ormDB = db
	initTable()
	return ormDB, nil
}

func DB() *gorm.DB {
	return ormDB
}

func CloseDB() {
	if ormDB != nil {
		sqlDB, err := ormDB.DB()
		if err != nil {
			return
		}
		sqlDB.Close()
	}
}

func initTable() {
	ormDB.AutoMigrate(&model.Member{})
}
