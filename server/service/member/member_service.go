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

func (s *MemberService) Auth(param *memberTypes.MemberAuthParam) (bool, error) {
	return true, nil
}

func (s *MemberService) AuthAndMember(param *memberTypes.MemberAuthParam) (*memberTypes.Member, error) {
	var errPreFix string = "Failed to auth member and get"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	db := s.DB
	username := param.Username
	var queryMember mysqlModel.Member
	if err := db.Where(mysqlModel.Member{
		Username: username,
	}).Take(&queryMember).Error; err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if !tool.CheckPassword(param.Password, queryMember.Password) {
		return nil, tool.PrefixError(errPreFix, errors.New("password is incorrect"))
	}
	return &memberTypes.Member{
		ID:       queryMember.ID,
		Username: queryMember.Username,
		Role:     memberTypes.MemberRole(queryMember.Role),
	}, nil
}

func (s *MemberService) Create(param *memberTypes.MemberCreateParam) error {
	var errPreFix string = "Failed to create member"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	// check user is existed
	db := s.DB
	var existedMember mysqlModel.Member
	result := db.Where(&mysqlModel.Member{Username: param.Username}).First(&existedMember)
	if result.Error == nil {
		return tool.PrefixError(errPreFix, errors.New("user is existed"))
	}

	// prepare member info for create
	pwd, err := tool.HashPassword(param.Password)
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	model := &mysqlModel.Member{
		Username: param.Username,
		Password: pwd,
		Email:    param.Email,
		Role:     int(memberTypes.Normal),
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

	fmt.Printf("member %s create successfully!", param.Username)
	return nil

}

func (s *MemberService) Edit(param *memberTypes.MemberEditParam) error {
	var errPreFix string = "Failed to member edit"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
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
	var errPreFix string = "Failed to member delete"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
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

func (s *MemberService) Member(param *memberTypes.MemberInfoParam) (*memberTypes.Member, error) {
	var errPreFix string = "Failed to get member"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	// find member
	db := s.DB
	queryMember := &mysqlModel.Member{
		Username: param.Username,
	}
	if err := db.Take(&queryMember).Error; err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	return &memberTypes.Member{
		ID:       queryMember.ID,
		Username: queryMember.Username,
		Role:     memberTypes.MemberRole(queryMember.Role),
	}, nil
}

func (s *MemberService) Members() ([]memberTypes.Member, error) {
	var errPreFix string = "Failed to get member list"
	var memberList []memberTypes.Member = make([]memberTypes.Member, 0)

	// check step
	err := s.CheckDB()
	if err != nil {
		return memberList, tool.PrefixError(errPreFix, err)
	}

	// find member list
	db := s.DB
	var memberModelList []mysqlModel.Member
	db.Find(&memberModelList)

	for _, memberModel := range memberModelList {
		memberList = append(memberList, memberTypes.Member{
			ID:       memberModel.ID,
			Username: memberModel.Username,
			Role:     memberTypes.MemberRole(memberModel.Role),
		})
	}

	return memberList, nil
}
