package types

import (
	"fmt"

	"gorm.io/gorm"
)

type Condition func(query *gorm.DB) *gorm.DB

func OrCondition(field, value string) Condition {
	return func(query *gorm.DB) *gorm.DB {
		return query.Or(
			fmt.Sprintf("%s LIKE ?", field),
			"%"+value+"%",
		)
	}
}

func AndCondition(field, value string) Condition {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(
			fmt.Sprintf("%s LIKE ?", field),
			"%"+value+"%",
		)
	}
}

func WithConditions(db *gorm.DB, conditions ...Condition) *gorm.DB {
	m := len(conditions)
	cond := make([]func(query *gorm.DB) *gorm.DB, 0, m)
	for _, condition := range conditions {
		cond = append(cond, condition)
	}
	return db.Scopes(cond...)
}
