package model

import "time"

type MemberProductLikes struct {
	MemberID  uint      `gorm:"primaryKey;index;uniqueIndex:member_product_likes_index"`
	ProductID uint      `gorm:"primaryKey;index;uniqueIndex:member_product_likes_index"`
	LikedAt   time.Time `gorm:"autoCreateTime"`
}

func (MemberProductLikes) TableName() string {
	return "member_product_likes"
}
