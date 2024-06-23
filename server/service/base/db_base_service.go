package base

import (
	"errors"

	"space.online.shop.web.server/service/db"
)

type DbBaseService struct {
	DB *db.DbService
}

func (srv *DbBaseService) SetDBService(DB *db.DbService) {
	srv.DB = DB
}

func (srv DbBaseService) CheckDB() error {
	if srv.DB.DB == nil {
		return errors.New("DB is nil")
	}
	return nil
}
