package base

import (
	"errors"

	"space.online.shop.web.server/service/db/mysql"
)

type DbBaseService struct {
	DB *mysql.MysqlService
}

func (srv *DbBaseService) SetDBService(DB *mysql.MysqlService) {
	srv.DB = DB
}

func (srv DbBaseService) CheckDB() error {
	if srv.DB.DB == nil {
		return errors.New("DB is nil")
	}
	return nil
}
