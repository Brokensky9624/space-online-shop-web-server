package product

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/base"
	"space.online.shop.web.server/service/db"
	mysqlModel "space.online.shop.web.server/service/db/model"
	dbTypes "space.online.shop.web.server/service/db/types"
	productTypes "space.online.shop.web.server/service/product/types"
	"space.online.shop.web.server/shared/utils/logger"
)

func NewService(DB *db.DbService) *ProductService {
	return &ProductService{
		DbBaseService: &base.DbBaseService{
			DB: DB,
		},
	}
}

type ProductService struct {
	*base.DbBaseService
}

func (s *ProductService) Create(userID uint, params ...productTypes.CreateParam) ([]uint, error) {
	idList := []uint{}

	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return idList, fmt.Errorf("database connection error")
	}

	idList, err := s.create(userID, params...)
	if err != nil {
		return idList, err
	}

	return idList, nil
}

func (s *ProductService) create(userID uint, params ...productTypes.CreateParam) ([]uint, error) {
	idList := []uint{}
	var anyError bool

	for _, param := range params {
		id, err := s.createOne(userID, param)
		if err != nil {
			anyError = true
			continue
		}
		idList = append(idList, id)
	}

	if anyError {
		return idList, fmt.Errorf("some product failed to create")
	}

	logger.SERVER.Info("succeed in create products, user_id: %v, product_id_list: %v", userID, idList)
	return idList, nil
}

func (s *ProductService) createOne(userID uint, param productTypes.CreateParam) (uint, error) {
	if err := param.Check(); err != nil {
		logger.SERVER.Error("invalid parameter, param: %+v, err: %v", userID, param, err)
		return 0, fmt.Errorf("invalid parameter")
	}

	createProduct := param.ToModel()
	if err := s.DB.Create(&createProduct).Error; err != nil {
		logger.SERVER.Error("failed to create product, user_id: %d, param: %+v, err: %v", userID, param, err)
		return 0, fmt.Errorf("failed to create product")
	}

	logger.SERVER.Info("succeed in creating one product, user_id: %v, product_id: %v", userID, createProduct.ID)
	return createProduct.ID, nil
}

func (s *ProductService) Edit(userID uint, param productTypes.EditParam) error {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return fmt.Errorf("database connection error")
	}

	if err := s.edit(userID, param); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) edit(userID uint, param productTypes.EditParam) error {
	if err := param.Check(); err != nil {
		logger.SERVER.Error("invalid parameter, user_id: %d, param: %+v, err: %v", userID, param, err)
		return fmt.Errorf("invalid parameter")
	}

	var queryProduct mysqlModel.Product
	if err := s.DB.First(&queryProduct, param.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.SERVER.Error("product not found, user_id: %v, product_id: %v", userID, param.ID)
			return fmt.Errorf("product not found")
		}
		logger.SERVER.Error("database error, user_id: %v, product_id: %v, err: %v", userID, param.ID, err)
		return fmt.Errorf("database error")
	}

	editProduct := param.ToModel()

	if err := s.DB.Model(&queryProduct).Updates(editProduct).Error; err != nil {
		logger.SERVER.Error("failed to updates product, user_id: %d, param: %+v, err: %v", userID, param, err)
		return fmt.Errorf("failed to updates product")
	}

	logger.SERVER.Info("succeed in updating product, user_id: %v, product_id: %v\n", userID, queryProduct.ID)
	return nil
}

func (s *ProductService) Like(userID uint, param productTypes.LikeParam) error {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return fmt.Errorf("database connection error")
	}

	if err := s.like(userID, param); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) like(userID uint, param productTypes.LikeParam) error {
	if err := param.Check(); err != nil {
		logger.SERVER.Error("invalid parameter, user_id: %d, param: %+v, err: %v", userID, param, err)
		return fmt.Errorf("invalid parameter")
	}

	var queryMember mysqlModel.Member
	if err := s.DB.First(&queryMember, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.SERVER.Error("member not found, user_id: %v", userID)
			return fmt.Errorf("member not found")
		}
		logger.SERVER.Error("database error, user_id: %v, err: %v", userID, err)
		return fmt.Errorf("database error")
	}

	var queryProduct mysqlModel.Product
	if err := s.DB.First(&queryProduct, param.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.SERVER.Error("product not found, user_id: %v, product_id: %v", userID, param.ProductID)
			return fmt.Errorf("product not found")
		}
		logger.SERVER.Error("database error, user_id: %v, product_id: %v, err: %v", userID, param.ProductID, err)
		return fmt.Errorf("database error")
	}

	if err := s.DB.Model(&queryMember).Association("Likes").Append(&queryProduct); err != nil {
		logger.SERVER.Error("failed to like product, user_id: %v, product_id: %v, err: %v", userID, queryProduct.ID, err)
		return fmt.Errorf("failed to like product")
	}

	logger.SERVER.Info("succeed in like product, user_id: %v, product_id: %v\n", userID, queryProduct.ID)
	return nil
}

func (s *ProductService) Delete(userID uint, param ...productTypes.DeleteParam) ([]uint, error) {
	var idList []uint

	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return idList, fmt.Errorf("database connection error")
	}

	if idList, err := s.delete(userID, param...); err != nil {
		return idList, err
	}
	return idList, nil
}

func (s *ProductService) delete(userID uint, params ...productTypes.DeleteParam) ([]uint, error) {
	var idList []uint
	var anyError bool

	for _, param := range params {
		id, err := s.deleteOne(userID, param)
		if err != nil {
			anyError = true
			continue
		}
		idList = append(idList, id)
	}

	if anyError {
		return idList, fmt.Errorf("some product failed to delete")
	}

	logger.SERVER.Info("succeed in delete products, user_id: %v, product_id_list: %v", userID, idList)
	return idList, nil
}

func (s *ProductService) deleteOne(userID uint, param productTypes.DeleteParam) (uint, error) {
	if err := param.Check(); err != nil {
		logger.SERVER.Error("invalid parameter, user_id: %d, param: %+v, err: %v", userID, param, err)
		return 0, fmt.Errorf("invalid parameter")
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		var queryProduct mysqlModel.Product
		if err := s.DB.First(&queryProduct, param.ProductID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.SERVER.Error("product not found, user_id: %v, product_id: %v", userID, param.ProductID)
				return fmt.Errorf("product not found")
			}
			logger.SERVER.Error("database error, user_id: %v, product_id: %v, err: %v", userID, param.ProductID, err)
			return fmt.Errorf("database error")
		}

		if err := s.DB.Model(&queryProduct).Association("LikedBy").Clear(); err != nil {
			logger.SERVER.Error("failed to clear likes, user_id: %v, product_id: %v, err: %v", userID, queryProduct.ID, err)
			return fmt.Errorf("failed to clear likes")
		}

		if err := s.DB.Delete(&queryProduct).Error; err != nil {
			logger.SERVER.Error("failed to delete product, user_id: %v, product_id: %v, err: %v", userID, queryProduct.ID, err)
			return fmt.Errorf("failed to delete product")
		}

		if err := s.DB.Unscoped().Delete(&queryProduct).Error; err != nil {
			logger.SERVER.Error("failed to delete unscoped product, user_id: %v, product_id: %v, err: %v", userID, queryProduct.ID, err)
			return fmt.Errorf("failed to delete unscoped product")
		}

		return nil
	}); err != nil {
		return 0, err
	}

	logger.SERVER.Info("succeed in deleting one product, user_id: %v, product_id: %v\n", userID, param.ProductID)
	return param.ProductID, nil
}

func (s *ProductService) Detail(param productTypes.DetailParam) (*productTypes.Product, error) {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return nil, fmt.Errorf("database connection error")
	}

	product, err := s.detail(param)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *ProductService) detail(param productTypes.DetailParam) (*productTypes.Product, error) {
	if err := param.Check(); err != nil {
		logger.SERVER.Error("invalid parameter, param: %+v, err: %v", param, err)
		return nil, fmt.Errorf("invalid parameter")
	}

	var queryProduct mysqlModel.Product
	if err := s.DB.First(&queryProduct, param.ProductID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.SERVER.Error("product not found, product_id: %v", param.ProductID)
			return nil, fmt.Errorf("product not found")
		}
		logger.SERVER.Error("database error, product_id: %v, err: %v", param.ProductID, err)
		return nil, fmt.Errorf("database error")
	}

	pd := productModelToProduct(queryProduct)

	return &pd, nil
}

func (s *ProductService) Query(
	conditions ...dbTypes.Condition,
) ([]productTypes.Product, error) {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return nil, fmt.Errorf("database connection error")
	}

	products, err := s.query(conditions...)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductService) query(
	conditions ...dbTypes.Condition,
) ([]productTypes.Product, error) {
	var queryProducts []mysqlModel.Product

	query := dbTypes.Scopes(
		s.DB.DB,
		conditions...,
	)

	if err := query.Find(&queryProducts).Error; err != nil {
		logger.SERVER.Error("database error, err: %v", err)
		return nil, fmt.Errorf("database error")
	}

	var products []productTypes.Product

	for _, p := range queryProducts {
		products = append(products, productModelToProduct(p))
	}

	return products, nil
}

func productModelToProduct(p mysqlModel.Product) productTypes.Product {
	return productTypes.Product{
		ID:           p.ID,
		Name:         p.Name,
		Title:        p.Title,
		Description:  p.Description,
		Category:     p.Category,
		Brand:        p.Brand,
		Manufacturer: p.Manufacturer,
		Status:       p.Status,
		Like:         uint(len(p.LikedBy)),
		UpdatedAt:    p.UpdatedAt,
		CreatedAt:    p.CreatedAt,
	}
}
