package types

import (
	dbTypes "space.online.shop.web.server/service/db/types"
)

type IProductSrv interface {
	Create(userID uint, params ...CreateParam) ([]uint, error)
	Edit(userID uint, param EditParam) error
	Like(userID uint, param LikeParam) error
	Delete(userID uint, param ...DeleteParam) ([]uint, error)
	Detail(param DetailParam) (*Product, error)
	Query(counter dbTypes.Counter, sortOrder dbTypes.SortOrder, searcher dbTypes.Searcher) ([]Product, error)
}
