package types

type IProductSrv interface {
	Create(userID uint, params ...CreateParam) ([]uint, error)
	Edit(userID uint, param EditParam) error
	Like(userID uint, param LikeParam) error
	Delete(userID uint, param ...DeleteParam) ([]uint, error)
	Detail(param DetailParam) (*Product, error)
	Query(qp QueryParam, op OrderParam) ([]Product, error)
}
