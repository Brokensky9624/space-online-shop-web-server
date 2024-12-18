package types

import "gorm.io/gorm"

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 10
)

type CounterOption interface {
	apply(c *counter)
}

type counterOptionFunc func(*counter)

func (fn counterOptionFunc) apply(c *counter) {
	fn(c)
}

func WithPage(page int) CounterOption {
	return counterOptionFunc(func(c *counter) {
		c.page = page
	})
}

func WithPageSize(pageSize int) CounterOption {
	return counterOptionFunc(func(c *counter) {
		c.pageSize = pageSize
	})
}

type Counter interface {
	Limit() int
	Offset() int
	Conditions() []Condition
}

type counter struct {
	conditions []Condition
	page       int
	pageSize   int
}

func (c *counter) Offset() int {
	if c.page < 1 || c.pageSize < 1 {
		return 0
	}
	return (c.page - 1) * c.pageSize
}

func (c *counter) Limit() int {
	if c.page < 1 || c.pageSize < 1 {
		return 0
	}
	return c.pageSize
}

func (c *counter) Conditions() []Condition {
	if c.conditions == nil {
		c.conditions = make([]Condition, 0)

		c.conditions = append(c.conditions, func(db *gorm.DB) *gorm.DB {
			return db.Offset(c.Offset()).Limit(c.Limit())
		})
	}

	return c.conditions
}

func NewCounter(counterOpts ...CounterOption) *counter {
	c := &counter{
		page:     DEFAULT_PAGE,
		pageSize: DEFAULT_PAGE_SIZE,
	}

	for _, opt := range counterOpts {
		opt.apply(c)
	}

	return c
}
