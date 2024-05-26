package product

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"space.online.shop.web.server/service/base"
	"space.online.shop.web.server/service/db/mysql"
	mysqlModel "space.online.shop.web.server/service/db/mysql/model"
	productTypes "space.online.shop.web.server/service/product/types"

	"space.online.shop.web.server/util/tool"
)

func NewService(DB *mysql.MysqlService) *ProdocutService {
	return &ProdocutService{
		DbBaseService: &base.DbBaseService{
			DB: DB,
		},
	}
}

type ProdocutService struct {
	*base.DbBaseService
}

// single
func (s *ProdocutService) Create(userID uint, param productTypes.CreateParam) error {
	var errPreFix string = "failed to create product"

	// check step
	if err := s.CheckDB(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	if err := param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	model := mysqlModel.ToProductModel(param)
	model.SetOwner(userID)

	if err := s.DB.Create(&model).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	return nil
}

func (s *ProdocutService) Edit(userID uint, param productTypes.EditParam) error {
	var errPreFix string = "failed to edit product"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	errPreFix = fmt.Sprintf("failed to edit product %d", param.ID)

	var queryModel mysqlModel.Product
	queryModel.SetID(param.ID)
	var matchModel mysqlModel.Product
	if err := s.DB.Where(queryModel).Take(&matchModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	if !matchModel.IsOwner(userID) {
		return tool.PrefixError(errPreFix, errors.New("you can not edit the product which does not belongs to you"))
	}

	editModel := mysqlModel.ToProductModel(param)
	if err := s.DB.Where(matchModel).Take(&matchModel).Updates(&editModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	fmt.Printf("product %d edit successfully!\n", editModel.ID)
	return nil
}

func (s *ProdocutService) Like(userID uint, param productTypes.LikeParam) error {
	var errPreFix string = "failed to like product"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	errPreFix = fmt.Sprintf("failed to like product %d", param.ID)

	var queryModel mysqlModel.Product
	queryModel.SetID(param.ID)
	var matchModel mysqlModel.Product
	editModel := mysqlModel.ToProductModel(param)
	if err := s.DB.Where(queryModel).Take(&matchModel).Updates(&editModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	fmt.Printf("member %d likes product %d successfully!\n", userID, editModel.ID)
	return nil
}

func (s *ProdocutService) Delete(userID uint, param productTypes.DeleteParam) error {
	var errPreFix string = "failed to delete product"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	errPreFix = fmt.Sprintf("failed to delete product %d", param.ID)

	var queryModel mysqlModel.Product
	queryModel.SetID(param.ID)
	var matchModel mysqlModel.Product

	if err := s.DB.Where(queryModel).Take(&matchModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	if !matchModel.IsOwner(userID) {
		return tool.PrefixError(errPreFix, errors.New("you can not edit the product which does not belongs to you"))
	}

	var deleteModel mysqlModel.Product
	deleteModel.SetID(matchModel.ID)
	if err := s.DB.Delete(&deleteModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	if err := s.DB.Unscoped().Delete(&deleteModel).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	fmt.Printf("member %d deletes product %d successfully!\n", userID, deleteModel.ID)
	return nil
}

func (s *ProdocutService) Detail(param productTypes.DetailParam) (*productTypes.Product, error) {
	var errPreFix string = "failed to get product detail"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	var model mysqlModel.Product
	model.ID = param.ID

	pd, err := s.queryProductByModel(model)
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	return pd, nil
}

// batches
func (s *ProdocutService) CreateInBatches(userID uint, params []productTypes.CreateParam) error {
	var errPreFix string = "failed to create product"

	// check step
	if err := s.CheckDB(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	for _, param := range params {
		if err := param.Check(); err != nil {
			return tool.PrefixError(errPreFix, err)
		}
	}
	models := []mysqlModel.Product{}
	for _, param := range params {
		model := mysqlModel.ToProductModel(param)
		model.OwnerID = userID
		models = append(models, model)
	}

	if err := s.DB.CreateInBatches(models, len(models)).Error; err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	return nil
}

func (s *ProdocutService) DeleteInBatches(userID uint, param productTypes.DeleteBatchesParam) error {
	var errPreFix string = "failed to delete products"

	// check step
	err := s.CheckDB()
	if err != nil {
		return tool.PrefixError(errPreFix, err)
	}
	if err = param.Check(); err != nil {
		return tool.PrefixError(errPreFix, err)
	}

	successIDList := []string{}

	var errSum error
	for _, ID := range param.IDList {
		errPreFix = fmt.Sprintf("failed to delete product %d", ID)

		var queryModel mysqlModel.Product
		queryModel.SetID(ID)
		var matchModel mysqlModel.Product

		if err := s.DB.Where(queryModel).Take(&matchModel).Error; err != nil {
			errSum = tool.MergeErrors(errSum, tool.PrefixError(errPreFix, err))
			continue
		}

		if !matchModel.IsOwner(userID) {
			errSum = tool.MergeErrors(errSum, errors.New("you can not edit the product which does not belongs to you"))
			continue
		}

		var deleteModel mysqlModel.Product
		deleteModel.SetID(matchModel.ID)
		if err := s.DB.Delete(&deleteModel).Error; err != nil {
			errSum = tool.MergeErrors(errSum, tool.PrefixError(errPreFix, err))
			continue
		}

		if err := s.DB.Unscoped().Delete(&deleteModel).Error; err != nil {
			errSum = tool.MergeErrors(errSum, tool.PrefixError(errPreFix, err))
			continue
		}

		successIDList = append(successIDList, strconv.FormatUint(uint64(ID), 10))
	}
	fmt.Printf("member %d deletes product (%s) successfully!\n", userID, strings.Join(successIDList, ","))
	return errSum
}

func (s *ProdocutService) Query() ([]productTypes.Product, error) {
	var errPreFix string = "failed to query product"

	// check step
	err := s.CheckDB()
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	return nil, nil
}

func (s *ProdocutService) queryProductByModel(model mysqlModel.Product) (*productTypes.Product, error) {
	if err := s.DB.Where(model).Take(&model).Error; err != nil {
		return nil, err
	}
	return ModelToProduct(model), nil
}

func ModelToProduct(m mysqlModel.Product) *productTypes.Product {
	product := productTypes.ToProduct(m)
	return &product
}
