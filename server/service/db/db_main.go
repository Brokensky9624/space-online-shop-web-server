package db

import (
	"context"
	"sync"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/db/interfaces"
)

var (
	once sync.Once
	db   *DbService
)

type DbService struct {
	*gorm.DB
	ctx          context.Context
	migratorList []interface{}
	wg           sync.WaitGroup
}

func (d *DbService) prepare() error {
	return d.AutoMigrate(d.migratorList...)
}

func (d *DbService) ping() error {
	// get sqlDB for close later
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *DbService) Close() error {
	// get sqlDB for close later
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *DbService) Run() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		<-d.ctx.Done()
		_ = d.Close()
	}()
}

func (d *DbService) Stop() {
	d.wg.Wait()
}

// = export functions
func NewDbService(ctx context.Context, dbBuilder interfaces.IDbBuilder) *DbService {
	once.Do(func() {
		gormDB, err := dbBuilder.BuildDB()
		if err != nil {
			panic(err)
		}

		tmpDb := &DbService{
			ctx: ctx,
			DB:  gormDB,
		}
		if err := tmpDb.ping(); err != nil {
			panic(err)
		}
		if err := tmpDb.prepare(); err != nil {
			panic(err)
		}
		db = tmpDb
	})
	return db
}

func Service() *DbService {
	return db
}
