package service

import (
	"space.online.shop.web.server/service/db/mysql"
	memberTypes "space.online.shop.web.server/service/member/types"
)

type IMemberSrv memberTypes.IMemberSrv

var manager *ServiceManager

func NewManager() *ServiceManager {
	manager = &ServiceManager{}
	return manager
}

func Manager() *ServiceManager {
	return manager
}

type ServiceManager struct {
	MemberSrv IMemberSrv
	MysqlSrv  *mysql.MysqlService
}

func (m *ServiceManager) SetMemberService(srv IMemberSrv) *ServiceManager {
	m.MemberSrv = srv
	return m
}

func (m *ServiceManager) SetDBService(srv *mysql.MysqlService) *ServiceManager {
	m.MysqlSrv = srv
	return m
}

func (m *ServiceManager) MemberService() IMemberSrv {
	return m.MemberSrv
}

func (m *ServiceManager) DBService() *mysql.MysqlService {
	return m.MysqlSrv
}
