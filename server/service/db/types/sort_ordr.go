package types

import (
	"fmt"

	"gorm.io/gorm"
)

type OrderType string

const (
	ASC  OrderType = "asc"
	DESC OrderType = "desc"
)

type SortOrderOption interface {
	apply(so *sortOrder)
}

type SortOrderOptionFunc func(*sortOrder)

func (fn SortOrderOptionFunc) apply(so *sortOrder) {
	fn(so)
}

func WithValidColumnsForSortOrder(validColumns map[string]struct{}) SortOrderOption {
	return SortOrderOptionFunc(func(so *sortOrder) {
		so.validColumns = validColumns
	})
}

func WithSortOrderMap(sortOrderMap map[string]string) SortOrderOption {
	return SortOrderOptionFunc(func(so *sortOrder) {
		so.sortOrderMap = sortOrderMap
	})
}

type SortOrder interface {
	Conditions() []Condition
}

type sortOrder struct {
	conditions   []Condition
	sortOrderMap map[string]string
	validColumns map[string]struct{}
}

func NewSortOrder(sortOrderOpts ...SortOrderOption) *sortOrder {
	so := &sortOrder{
		sortOrderMap: make(map[string]string),
		validColumns: make(map[string]struct{}),
	}

	for _, opt := range sortOrderOpts {
		opt.apply(so)
	}

	return so
}

func (s *sortOrder) Conditions() []Condition {
	if s.conditions == nil {
		s.conditions = make([]Condition, 0)

		for column, order := range s.sortOrderMap {
			if _, ok := s.validColumns[column]; !ok {
				continue
			}

			switch order {
			case string(ASC), string(DESC):
				s.conditions = append(s.conditions, func(db *gorm.DB) *gorm.DB {
					return db.Order(fmt.Sprintf("%s %s", column, order))
				})
			}
		}
	}

	return s.conditions
}
