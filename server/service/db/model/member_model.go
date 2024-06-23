package model

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Account  string `gorm:"size:100;not null"`
	Username string `gorm:"size:50;not null"`
	Password string `gorm:"size:200;not null"`
	Role     int    `gorm:"size:50;not null"`
	Email    string `gorm:"size:200;not null"`
	Phone    string `gorm:"size:100;not null"`
	Address  string `gorm:"size:200;not null"`
}

func (Member) TableName() string {
	return "member"
}
