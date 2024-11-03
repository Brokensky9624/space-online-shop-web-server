package product

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/base"
	"space.online.shop.web.server/service/db"
	mysqlModel "space.online.shop.web.server/service/db/model"
	productTypes "space.online.shop.web.server/service/product/types"

	"space.online.shop.web.server/shared/utils/logger"
	"space.online.shop.web.server/shared/utils/tool"
)

func NewService(DB *db.DbService) *ProdocutService {
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
		return tool.PrefixError(fmt.Sprintf("%s: database connection error, user_id: %v", errPreFix, userID), err)
	}

	if err := param.Check(); err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: parameter error, user_id: %v", errPreFix, userID), err)
	}

	createProduct := param.ToModel()

	if err := s.DB.Create(&createProduct).Error; err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: create error, user_id: %v, product_id: %v", errPreFix, userID, createProduct.ID), err)
	}

	logger.SERVER.Info(fmt.Sprintf("succeed to create product, user_id: %v, product_id: %v", userID, createProduct.ID))

	return nil
}

func (s *ProdocutService) Edit(userID uint, param productTypes.EditParam) error {
	var errPreFix string = "failed to edit product"

	// check step
	if err := s.CheckDB(); err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: database connection error, user_id: %v", errPreFix, userID), err)
	}

	if err := param.Check(); err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: parameter error, user_id: %v", errPreFix, userID), err)
	}

	var queryProduct mysqlModel.Product
	if err := s.DB.First(&queryProduct, param.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tool.PrefixError(fmt.Sprintf("%s: product not found, user_id: %v, param_id: %v", errPreFix, userID, param.ID), err)
		}
		return tool.PrefixError(fmt.Sprintf("%s: database error, user_id: %v, param_id: %v", errPreFix, userID, param.ID), err)
	}

	editProduct := param.ToModel()

	if err := s.DB.Model(&queryProduct).Updates(editProduct).Error; err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: updates error, user_id: %v, product_id: %v", errPreFix, userID, queryProduct.ID), err)
	}

	logger.SERVER.Info("succeed to edit product, user_id: %v, product_id: %v\n", userID, queryProduct.ID)
	return nil
}

func (s *ProdocutService) Like(userID uint, param productTypes.LikeParam) error {
	var errPreFix string = "failed to like product"

	// check step
	if err := s.CheckDB(); err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: database connection error, user_id: %v", errPreFix, userID), err)
	}

	if err := param.Check(); err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: parameter error, user_id: %v", errPreFix, userID), err)
	}

	var queryMember mysqlModel.Member
	if err := s.DB.First(&queryMember, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tool.PrefixError(fmt.Sprintf("%s: member not found, user_id: %v", errPreFix, userID), err)
		}
		return tool.PrefixError(fmt.Sprintf("%s: database error, user_id: %v", errPreFix, userID), err)
	}

	var queryProduct mysqlModel.Product
	if err := s.DB.First(&queryProduct, param.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tool.PrefixError(fmt.Sprintf("%s: product not found, user_id: %v, param_id: %v", errPreFix, userID, param.ProductID), err)
		}
		return tool.PrefixError(fmt.Sprintf("%s: database error, user_id: %v, param_id: %v", errPreFix, userID, param.ProductID), err)
	}

	if err := s.DB.Model(&queryProduct).Association("Likes").Append(&queryMember); err != nil {
		return tool.PrefixError(fmt.Sprintf("%s: append error, user_id: %v, product_id: %v", errPreFix, userID, queryProduct.ID), err)
	}

	logger.SERVER.Info("succeed to like product, user_id: %v, product_id: %v\n", userID, queryProduct.ID)
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
