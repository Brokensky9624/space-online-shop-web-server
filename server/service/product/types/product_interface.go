package types

type IProductSrv interface {
	Create(userID uint, param CreateParam) error
	CreateInBatches(userID uint, params []CreateParam) error
	Edit(userID uint, param EditParam) error
	Like(userID uint, param LikeParam) error
	Delete(userID uint, param DeleteParam) error
	DeleteInBatches(userID uint, param DeleteBatchesParam) error
	Detail(param DetailParam) (*Product, error)
	Query() ([]Product, error)
}
