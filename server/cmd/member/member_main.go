package member

import (
	"errors"
	"fmt"

	"space.online.shop.web.server/service/db/model"
	mysqlSrv "space.online.shop.web.server/service/db/mysql"
	"space.online.shop.web.server/util/tool"
)

type param = map[string]interface{}

func Create(memberParam param) error {
	var errPreFix string = "failed to member create"
	db := mysqlSrv.DB()
	if db == nil {
		return tool.PrefixError(errPreFix, errors.New("mysql db is nil"))
	}
	username, ok := memberParam["username"].(string)
	if !ok || username == "" {
		return tool.PrefixError(errPreFix, errors.New("need Username"))
	}
	password, ok := memberParam["password"].(string)
	if !ok || password == "" {
		return tool.PrefixError(errPreFix, errors.New("need Password"))
	}
	pwd, err := tool.HashPassword(password)
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	email, _ := memberParam["email"].(string)
	model := &model.Member{
		Username: username,
		Password: pwd,
		Email:    email,
	}
	result := db.Create(model)
	if result.Error != nil {
		return tool.PrefixError(errPreFix, result.Error)
	}
	fmt.Printf("member %s create successfully!", username)
	return nil
}

func Edit(memberParam param) error {
	var errPreFix string = "failed to member edit"
	db := mysqlSrv.DB()
	if db == nil {
		return tool.PrefixError(errPreFix, errors.New("mysql db is nil"))
	}
	username, ok := memberParam["username"].(string)
	if !ok || username == "" {
		return tool.PrefixError(errPreFix, errors.New("need username"))
	}
	email, ok := memberParam["email"].(string)
	if !ok || email == "" {
		return tool.PrefixError(errPreFix, errors.New("need email"))
	}
	memberNewInfo := &model.Member{
		Email: email,
	}
	result := db.Model(&model.Member{}).Where("username = ?", username).Updates(memberNewInfo)
	if result.Error != nil {
		return tool.PrefixError(errPreFix, result.Error)
	}
	if result.RowsAffected == 0 {
		return tool.PrefixError(errPreFix, errors.New("no member update"))
	}
	fmt.Printf("member %s edit successfully!", username)
	return nil
}

func Delete() error {
	var err error
	var errPreFix string = "failed to member delete"
	defer func() {
		if err != nil {
			err = tool.PrefixError(errPreFix, err)
		}
	}()
	db := mysqlSrv.DB()
	if db == nil {
		return errors.New("mysql db is nil")
	}
	return err
}

func List() ([]model.Member, error) {
	var err error
	var errPreFix string = "failed to member list"
	defer func() {
		if err != nil {
			err = tool.PrefixError(errPreFix, err)
		}
	}()
	db := mysqlSrv.DB()
	if db == nil {
		return nil, errors.New("mysql db is nil")
	}
	var members []model.Member
	db.Find(&members)
	return members, err
}
