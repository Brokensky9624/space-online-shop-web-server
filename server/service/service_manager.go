package service

import (
	"space.online.shop.web.server/service/db/mysql"
	memberTypes "space.online.shop.web.server/service/member/types"
	productTypes "space.online.shop.web.server/service/product/types"
)

type IMemberSrv memberTypes.IMemberSrv
type IProductSrv productTypes.IProductSrv

var manager *ServiceManager

func NewManager() *ServiceManager {
	manager = &ServiceManager{}
	return manager
}

func Manager() *ServiceManager {
	return manager
}

type ServiceManager struct {
	MysqlSrv   *mysql.MysqlService
	MemberSrv  IMemberSrv
	ProductSrv IProductSrv
}

func (m *ServiceManager) SetDBService(srv *mysql.MysqlService) *ServiceManager {
	m.MysqlSrv = srv
	return m
}

func (m *ServiceManager) SetMemberService(srv IMemberSrv) *ServiceManager {
	m.MemberSrv = srv
	return m
}

func (m *ServiceManager) SetProductService(srv IProductSrv) *ServiceManager {
	m.ProductSrv = srv
	return m
}

func (m *ServiceManager) DBService() *mysql.MysqlService {
	return m.MysqlSrv
}

func (m *ServiceManager) MemberService() IMemberSrv {
	return m.MemberSrv
}

func (m *ServiceManager) ProductService() IProductSrv {
	return m.ProductSrv
}
