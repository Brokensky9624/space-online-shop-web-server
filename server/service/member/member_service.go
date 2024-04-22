package member

import (
	"errors"
	"fmt"

	"space.online.shop.web.server/service/db/mysql"
	mysqlModel "space.online.shop.web.server/service/db/mysql/model"
	memberTypes "space.online.shop.web.server/service/member/types"

	"space.online.shop.web.server/util/tool"
)

func NewService() *MemberService {
	return &MemberService{}
}

type MemberService struct {
	DB *mysql.MysqlService
}

func (s *MemberService) SetDBService(DB *mysql.MysqlService) *MemberService {
	s.DB = DB
	return s
}

func (s *MemberService) CheckDB() error {
	if s.DB.DB == nil {
		return fmt.Errorf("DB of Member service is nil")
	}
	return nil
}

func (s *MemberService) Auth() (bool, error) {
	return true, nil
}

func (s *MemberService) AuthAndMember() (*memberTypes.Member, error) {
	return nil, nil
}

func (s *MemberService) Create(param *memberTypes.MemberCreateParam) error {
	var errPreFix string = "failed to member create"

	// check step
	err := s.CheckDB()
	if err != nil {
		return err
	}
	if err = param.Check(); err != nil {
		return err
	}

	// check user is existed
	db := s.DB
	username := param.Username
	var existedMember mysqlModel.Member
	result := db.Where(&mysqlModel.Member{Username: username}).First(&existedMember)
	if result.Error == nil {
		return tool.PrefixError(errPreFix, errors.New("user is existed"))
	}

	// prepare member info for create
	pwd, err := tool.HashPassword(param.Password)
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	email := param.Email
	model := &mysqlModel.Member{
		Username: username,
		Password: pwd,
		Email:    email,
	}

	// create member
	result = db.Create(model)
	if result.Error != nil {
		return tool.PrefixError(errPreFix, result.Error)
	}

	// check affect
	if result.RowsAffected == 0 {
		return tool.PrefixError(errPreFix, errors.New("no member created"))
	}

	fmt.Printf("member %s create successfully!", username)
	return nil

}

func (s *MemberService) Edit(param *memberTypes.MemberEditParam) error {
	var errPreFix string = "failed to member edit"

	// check step
	err := s.CheckDB()
	if err != nil {
		return err
	}
	if err = param.Check(); err != nil {
		return err
	}

	// update member
	db := s.DB
	username := param.Username
	memberNewInfo := &mysqlModel.Member{
		Email: param.Email,
	}
	result := db.Model(&mysqlModel.Member{}).Where("username = ?", username).Updates(memberNewInfo)
	if result.Error != nil {
		return tool.PrefixError(errPreFix, result.Error)
	}

	// check affect
	if result.RowsAffected == 0 {
		return tool.PrefixError(errPreFix, errors.New("no member updated"))
	}
	fmt.Printf("member %s edit successfully!", username)
	return nil
}

func (s *MemberService) Delete(param *memberTypes.MemberDeleteParam) error {
	var errPreFix string = "failed to member delete"

	// check step
	err := s.CheckDB()
	if err != nil {
		return err
	}
	if err = param.Check(); err != nil {
		return err
	}

	// delete member
	db := s.DB
	username := param.Username
	result := db.Where("username = ?", username).Delete(&mysqlModel.Member{})
	if result.Error != nil {
		return tool.PrefixError(errPreFix, result.Error)
	}
	// check affect
	if result.RowsAffected == 0 {
		return tool.PrefixError(errPreFix, errors.New("no member deleted"))
	}
	fmt.Printf("member %s delete successfully!", username)
	return nil
}

func (s *MemberService) Member() (*memberTypes.Member, error) {
	return nil, nil
}

func (s *MemberService) Members() ([]memberTypes.Member, error) {
	var errPreFix string = "failed to member list"
	var memberList []memberTypes.Member = make([]memberTypes.Member, 0)

	// check step
	err := s.CheckDB()
	if err != nil {
		return memberList, tool.PrefixError(errPreFix, err)
	}

	// find member list
	db := s.DB
	db.Find(&memberList)

	// 	return members, nil
	return memberList, nil
}
