package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"space.online.shop.web.server/service/db/mysql/model"
)

var db *MysqlService

func Service() *MysqlService {
	return db
}

func NewDBService() (*MysqlService, error) {
	// FIXME: load cfg information by env
	gormDB, err := NewDBBuilder().
		SetIp("localhost").
		SetPort("3306").
		SetUserName("space_online_admin").
		SetPassword("space_online_is_666").
		SetDatabaseName("masterDB").
		Build()
	if err != nil {
		return nil, err
	}
	db = &MysqlService{DB: gormDB}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err = db.InitTable(); err != nil {
		return nil, err
	}
	return db, nil
}

func NewDBBuilder() *DBBuilder {
	return &DBBuilder{
		cfg: &dbCfg{},
	}
}

type DBBuilder struct {
	cfg *dbCfg
}

type dbCfg struct {
	ip           string
	port         string
	userName     string
	password     string
	databaseName string
}

func (b *DBBuilder) SetIp(ip string) *DBBuilder {
	b.cfg.ip = ip
	return b
}

func (b *DBBuilder) SetPort(port string) *DBBuilder {
	b.cfg.port = port
	return b
}

func (b *DBBuilder) SetUserName(userName string) *DBBuilder {
	b.cfg.userName = userName
	return b
}

func (b *DBBuilder) SetPassword(password string) *DBBuilder {
	b.cfg.password = password
	return b
}

func (b *DBBuilder) SetDatabaseName(databaseName string) *DBBuilder {
	b.cfg.databaseName = databaseName
	return b
}

func (b *DBBuilder) Build() (*gorm.DB, error) {
	cfg := b.cfg
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.userName, cfg.password, cfg.ip, cfg.port, cfg.databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type MysqlService struct {
	*gorm.DB
}

func (d *MysqlService) InitTable() error {
	return d.AutoMigrate(
		&model.Member{},
		&model.Product{},
	)
}

func (d *MysqlService) Ping() error {
	// get sqlDB for close later
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *MysqlService) Close() error {
	// get sqlDB for close later
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
