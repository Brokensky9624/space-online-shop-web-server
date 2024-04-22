package member

// import (
// 	"errors"
// 	"fmt"

// 	"space.online.shop.web.server/service/db/model"
// 	mysqlSrv "space.online.shop.web.server/service/db/mysql"
// 	"space.online.shop.web.server/util/tool"
// )

// type param = map[string]interface{}

// func Register(memberParam param) error {
// 	var errPreFix string = "failed to member create"
// 	db := mysqlSrv.DB()
// 	// check db is nil
// 	if db == nil {
// 		return tool.PrefixError(errPreFix, errors.New("mysql db is nil"))
// 	}
// 	// check username is empty
// 	username, ok := memberParam["username"].(string)
// 	if !ok || username == "" {
// 		return tool.PrefixError(errPreFix, errors.New("need username"))
// 	}
// 	// check password is empty
// 	password, ok := memberParam["password"].(string)
// 	if !ok || password == "" {
// 		return tool.PrefixError(errPreFix, errors.New("need password"))
// 	}
// 	// check user is existed
// 	var existedMember model.Member
// 	result := db.Where(&model.Member{Username: username}).First(&existedMember)
// 	if result.Error == nil {
// 		return tool.PrefixError(errPreFix, errors.New("user is existed"))
// 	}
// 	// has password
// 	pwd, err := tool.HashPassword(password)
// 	if err != nil {
// 		return tool.PrefixError(errPreFix, err)
// 	}
// 	email, _ := memberParam["email"].(string)
// 	model := &model.Member{
// 		Username: username,
// 		Password: pwd,
// 		Email:    email,
// 	}
// 	// create member
// 	result = db.Create(model)
// 	if result.Error != nil {
// 		return tool.PrefixError(errPreFix, result.Error)
// 	}
// 	// check affect
// 	if result.RowsAffected == 0 {
// 		return tool.PrefixError(errPreFix, errors.New("no member created"))
// 	}
// 	fmt.Printf("member %s create successfully!", username)
// 	return nil
// }

// func Edit(memberParam param) error {
// 	var errPreFix string = "failed to member edit"
// 	// check db is nil
// 	db := mysqlSrv.DB()
// 	if db == nil {
// 		return tool.PrefixError(errPreFix, errors.New("mysql db is nil"))
// 	}
// 	// check username is empty
// 	username, ok := memberParam["username"].(string)
// 	if !ok || username == "" {
// 		return tool.PrefixError(errPreFix, errors.New("need username"))
// 	}
// 	// check email is empty
// 	email, ok := memberParam["email"].(string)
// 	if !ok || email == "" {
// 		return tool.PrefixError(errPreFix, errors.New("need email"))
// 	}
// 	memberNewInfo := &model.Member{
// 		Email: email,
// 	}
// 	// update member
// 	result := db.Model(&model.Member{}).Where("username = ?", username).Updates(memberNewInfo)
// 	if result.Error != nil {
// 		return tool.PrefixError(errPreFix, result.Error)
// 	}
// 	// check affect
// 	if result.RowsAffected == 0 {
// 		return tool.PrefixError(errPreFix, errors.New("no member updated"))
// 	}
// 	fmt.Printf("member %s edit successfully!", username)
// 	return nil
// }

// func Delete(username string) error {
// 	var errPreFix string = "failed to member delete"
// 	db := mysqlSrv.DB()
// 	// check db is nil
// 	if db == nil {
// 		return tool.PrefixError(errPreFix, errors.New("mysql db is nil"))
// 	}
// 	// check username is empty
// 	if username == "" {
// 		return tool.PrefixError(errPreFix, errors.New("need username"))
// 	}
// 	// delete member
// 	result := db.Where("username = ?", username).Delete(&model.Member{})
// 	if result.Error != nil {
// 		return tool.PrefixError(errPreFix, result.Error)
// 	}
// 	// check affect
// 	if result.RowsAffected == 0 {
// 		return tool.PrefixError(errPreFix, errors.New("no member deleted"))
// 	}
// 	fmt.Printf("member %s delete successfully!", username)
// 	return nil
// }

// func List() ([]model.Member, error) {
// 	var errPreFix string = "failed to member list"
// 	db := mysqlSrv.DB()
// 	// check db is nil
// 	if db == nil {
// 		return nil, tool.PrefixError(errPreFix, errors.New("mysql db is nil"))
// 	}
// 	// find all member
// 	var members []model.Member
// 	db.Find(&members)
// 	return members, nil
// }
