package member

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/base"
	"space.online.shop.web.server/service/db"
	mysqlModel "space.online.shop.web.server/service/db/model"
	memberTypes "space.online.shop.web.server/service/member/types"

	"space.online.shop.web.server/util/logger"
	"space.online.shop.web.server/util/tool"
)

func NewService(DB *db.DbService) *MemberService {
	return &MemberService{
		DbBaseService: &base.DbBaseService{
			DB: DB,
		},
	}
}

type MemberService struct {
	*base.DbBaseService
}

func (s *MemberService) Auth(param *memberTypes.MemberAuthParam) error {
	var errPreFix string = "failed to auth member"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	model := mysqlModel.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(model, true)
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if !tool.CheckPassword(param.Password, matchMember.Password) {
		return tool.PrefixError(errPreFix, errors.New("password is incorrect"))
	}
	fmt.Printf("member %s auth successfully!\n", param.Account)
	return nil
}

func (s *MemberService) AuthAndMember(param *memberTypes.MemberAuthParam) (*memberTypes.Member, error) {
	var errPreFix string = "failed to auth member and get"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	model := mysqlModel.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(model, true)
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if !tool.CheckPassword(param.Password, matchMember.Password) {
		return nil, tool.PrefixError(errPreFix, errors.New("password is incorrect"))
	}
	fmt.Printf("member %s auth successfully!\n", param.Account)
	return matchMember, nil
}

func (s *MemberService) Create(param memberTypes.MemberCreateParam) error {
	var errPreFix string = "failed to create member"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	// check user is existed
	model := mysqlModel.Member{
		Account: param.Account,
	}
	if _, err = s.queryMemberByModel(model, false); err == nil {
		return tool.PrefixError(errPreFix, errors.New("user is existed"))
	}

	// prepare member info for create
	pwd, err := tool.HashPassword(param.Password)
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	createModel := &mysqlModel.Member{
		Account:  param.Account,
		Username: param.Username,
		Password: pwd,
		Role:     int(memberTypes.Normal),
		Email:    param.Email,
		Phone:    param.Phone,
		Address:  param.Address,
	}

	// create member
	if err = s.DB.Create(createModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	fmt.Printf("member %s create successfully!\n", param.Account)
	return nil
}

func (s *MemberService) Edit(param memberTypes.MemberEditParam) error {
	var errPreFix string = "failed to member edit"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	// update member
	model := mysqlModel.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(model, false)
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	queryModel := mysqlModel.Member{
		Model: gorm.Model{
			ID: matchMember.ID,
		},
	}
	editModel := mysqlModel.Member{
		Account:  param.Account,
		Username: param.Username,
		Email:    param.Email,
		Phone:    param.Phone,
		Address:  param.Address,
	}
	if err := s.DB.Where(queryModel).Take(&queryModel).Updates(&editModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	fmt.Printf("member %s edit successfully!\n", param.Account)
	return nil
}

func (s *MemberService) Delete(param memberTypes.MemberDeleteParam) error {
	var errPreFix string = "failed to member delete"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	// delete member
	account := param.Account
	deleteMember := mysqlModel.Member{
		Account: account,
	}
	if err := s.DB.Where(deleteMember).Take(&deleteMember).Delete(&deleteMember).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	deleteUnscopedMember := mysqlModel.Member{
		Account: account,
	}
	if err := s.DB.Unscoped().Where(deleteUnscopedMember).Take(&deleteUnscopedMember).Delete(&deleteUnscopedMember).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	fmt.Printf("member %s delete successfully!\n", param.Account)
	return nil
}

func (s *MemberService) Member(param memberTypes.MemberInfoParam) (*memberTypes.Member, error) {
	var errPreFix string = "failed to get member"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	// find member list
	model := mysqlModel.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(model, false)
	return matchMember, tool.PrefixError(errPreFix, err)
}

func (s *MemberService) Members() ([]memberTypes.Member, error) {
	var errPreFix string = "failed to get member list"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	logger.SERVER.Debug("Test debug")
	logger.SERVER.Info("Test info")
	// find member list
	var memberList []memberTypes.Member = make([]memberTypes.Member, 0)
	var queryMemberList []mysqlModel.Member
	if err := s.DB.Find(&queryMemberList).Error; err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	for _, memberModel := range queryMemberList {
		memberList = append(memberList, *memberTypes.ModelToMember(memberModel, false))
	}
	return memberList, nil
}

func (s *MemberService) queryMemberByModel(model mysqlModel.Member, includePassword bool) (*memberTypes.Member, error) {
	var queryMemberList []mysqlModel.Member
	if err := s.DB.Where(model).Find(&queryMemberList).Error; err != nil {
		return nil, err
	}
	var matchMember *memberTypes.Member
	for _, queryMember := range queryMemberList {
		if queryMember.Account == model.Account {
			matchMember = memberTypes.ModelToMember(queryMember, includePassword)
			break
		}
	}
	if matchMember == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return matchMember, nil
}
