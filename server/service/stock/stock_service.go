package stock

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"space.online.shop.web.server/service/base"
	"space.online.shop.web.server/service/db"
	"space.online.shop.web.server/shared/utils/logger"
	"space.online.shop.web.server/shared/utils/tool"

	mysqlModel "space.online.shop.web.server/service/db/model"
	dbTypes "space.online.shop.web.server/service/db/types"
	stockTypes "space.online.shop.web.server/service/stock/types"
)

type StockService struct {
	*base.DbBaseService
}

func NewStockService(DB *db.DbService) *StockService {
	return &StockService{
		DbBaseService: &base.DbBaseService{
			DB: DB,
		},
	}
}

func (s *StockService) Create(userID uint, params ...stockTypes.CreateParam) ([]uint, error) {
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

func (s *StockService) create(userID uint, params ...stockTypes.CreateParam) ([]uint, error) {
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
		return idList, fmt.Errorf("some stock failed to create")
	}

	logger.SERVER.Info("succeed in create stocks, user_id: %v, stock_id_list: %v", userID, idList)
	return idList, nil
}

func (s *StockService) createOne(userID uint, param stockTypes.CreateParam) (uint, error) {
	if err := tool.CheckRequiredFields(param); err != nil {
		logger.SERVER.Error("invalid parameter, user_id: %d, param: %+v, err: %v", userID, param, err)
		return 0, fmt.Errorf("invalid parameter")
	}

	createStock := param.ToModel()

	if err := s.DB.Create(&createStock).Error; err != nil {
		logger.SERVER.Error("failed to create stock, user_id: %d, param: %+v, err: %v", userID, param, err)
		return 0, fmt.Errorf("failed to create stock")
	}

	logger.SERVER.Info("succeed in creating one stock, user_id: %v, stock_id: %v", userID, createStock.ID)
	return createStock.ID, nil
}

func (s *StockService) Edit(userID uint, param stockTypes.EditParam) error {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return fmt.Errorf("database connection error")
	}

	if err := s.edit(userID, param); err != nil {
		return err
	}

	return nil
}

func (s *StockService) edit(userID uint, param stockTypes.EditParam) error {
	if err := tool.CheckRequiredFields(param); err != nil {
		logger.SERVER.Error("invalid parameter, user_id: %d, param: %+v, err: %v", userID, param, err)
		return fmt.Errorf("invalid parameter")
	}

	var queryStock mysqlModel.Stock
	if err := s.DB.First(&queryStock, param.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.SERVER.Error("stock not found, user_id: %v, stock_id: %v", userID, param.ID)
			return fmt.Errorf("stock not found")
		}
		logger.SERVER.Error("database error, user_id: %v, stock_id: %v, err: %v", userID, param.ID, err)
		return fmt.Errorf("database error")
	}

	editStock := param.ToModel()

	if err := s.DB.Model(&queryStock).Updates(editStock).Error; err != nil {
		logger.SERVER.Error("failed to updates stock, user_id: %d, param: %+v, err: %v", userID, param, err)
		return fmt.Errorf("failed to updates stock")
	}

	logger.SERVER.Info("succeed in updating stock, user_id: %v, stock_id: %v\n", userID, queryStock.ID)
	return nil
}

func (s *StockService) Delete(userID uint, param ...stockTypes.DeleteParam) ([]uint, error) {
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

func (s *StockService) delete(userID uint, param ...stockTypes.DeleteParam) ([]uint, error) {
	var idList []uint
	var anyError bool

	for _, param := range param {
		id, err := s.deleteOne(userID, param)
		if err != nil {
			anyError = true
			continue
		}
		idList = append(idList, id)
	}

	if anyError {
		return idList, fmt.Errorf("some stock failed to delete")
	}

	logger.SERVER.Info("succeed in delete stocks, user_id: %v, stock_id_list: %v", userID, idList)
	return idList, nil
}

func (s *StockService) deleteOne(userID uint, param stockTypes.DeleteParam) (uint, error) {
	if err := tool.CheckRequiredFields(param); err != nil {
		logger.SERVER.Error("invalid parameter, user_id: %d, param: %+v, err: %v", userID, param, err)
		return 0, fmt.Errorf("invalid parameter")
	}

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		var queryStock mysqlModel.Stock
		if err := s.DB.First(&queryStock, param.StockID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.SERVER.Error("stock not found, user_id: %v, stock_id: %v", userID, param.StockID)
				return fmt.Errorf("stock not found")
			}
			logger.SERVER.Error("database error, user_id: %v, stock_id: %v, err: %v", userID, param.StockID, err)
			return fmt.Errorf("database error")
		}

		if err := s.DB.Delete(&queryStock).Error; err != nil {
			logger.SERVER.Error("failed to delete stock, user_id: %v, stock_id: %v, err: %v", userID, queryStock.ID, err)
			return fmt.Errorf("failed to delete stock")
		}

		if err := s.DB.Unscoped().Delete(&queryStock).Error; err != nil {
			logger.SERVER.Error("failed to delete unscoped stock, user_id: %v, stock_id: %v, err: %v", userID, queryStock.ID, err)
			return fmt.Errorf("failed to delete unscoped stock")
		}

		return nil
	}); err != nil {
		return 0, err
	}

	logger.SERVER.Info("succeed in deleting one stock, user_id: %v, stock_id: %v", userID, param.StockID)
	return param.StockID, nil
}

func (s *StockService) Detail(parm stockTypes.DetailParam) (*stockTypes.Stock, error) {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return nil, fmt.Errorf("database connection error")
	}

	return s.detail(parm)
}

func (s *StockService) detail(param stockTypes.DetailParam) (*stockTypes.Stock, error) {
	if err := tool.CheckRequiredFields(param); err != nil {
		logger.SERVER.Error("invalid parameter, param: %+v, err: %v", param, err)
		return nil, fmt.Errorf("invalid parameter")
	}

	var queryStock mysqlModel.Stock

	if err := s.DB.First(&queryStock, param.StockID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.SERVER.Error("stock not found, stock_id: %v", param.StockID)
			return nil, fmt.Errorf("stock not found")
		}
		logger.SERVER.Error("database error, stock_id: %v, err: %v", param.StockID, err)
		return nil, fmt.Errorf("database error")
	}

	stock := stockModelToStock(queryStock)
	return &stock, nil
}

func (s *StockService) Query(conditions ...dbTypes.Condition) ([]stockTypes.Stock, error) {
	if err := s.CheckDB(); err != nil {
		logger.SERVER.Error("database connection error, err: %v", err)
		return nil, fmt.Errorf("database connection error")
	}

	return s.query(conditions...)
}

func (s *StockService) query(conditions ...dbTypes.Condition) ([]stockTypes.Stock, error) {
	var queryStocks []mysqlModel.Stock

	query := dbTypes.Scopes(
		s.DB.DB,
		conditions...,
	)

	if err := query.Find(&queryStocks).Error; err != nil {
		logger.SERVER.Error("database error, err: %v", err)
		return nil, fmt.Errorf("database error")
	}

	var stocks []stockTypes.Stock

	for _, s := range queryStocks {
		stocks = append(stocks, stockModelToStock(s))
	}

	return stocks, nil
}

func stockModelToStock(s mysqlModel.Stock) stockTypes.Stock {
	return stockTypes.Stock{
		ID:          s.ID,
		ProductID:   s.ProductID,
		Quantity:    s.Quantity,
		Price:       s.Price,
		Description: s.Description,
		UpdatedAt:   s.UpdatedAt,
		CreatedAt:   s.CreatedAt,
	}
}
