package types

import (
	dbTypes "space.online.shop.web.server/service/db/types"
)

type IStockSrv interface {
	Create(userID uint, params ...CreateParam) ([]uint, error)
	Edit(userID uint, param EditParam) error
	Delete(userID uint, param ...DeleteParam) ([]uint, error)
	Detail(param DetailParam) (*Stock, error)
	Query(conditions ...dbTypes.Condition) ([]Stock, error)
}
