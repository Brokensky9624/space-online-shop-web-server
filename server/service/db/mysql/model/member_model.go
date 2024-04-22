package model

import "gorm.io/gorm"

type Member struct {
	gorm.Model
	Username string `gorm:"size:50;not null"`
	Password string `gorm:"size:200;not null"`
	Email    string `gorm:"size:200;not null"`
	Role     string `gorm:"size:50;not null"`
}

func (Member) TableName() string {
	return "member"
}
