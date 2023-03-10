package model

import (
	"time"

	"github.com/google/uuid"
)

// UserModel 用户模型
type UserModel struct {
	Id               uuid.UUID `gorm:"column:id" json:"id"`
	OpenId           string    `gorm:"column:open_id" json:"openId"`
	AvatarUrl        string    `gorm:"column:avatar_url" json:"avatarUrl"`
	City             string    `gorm:"column:city" json:"city"`
	Country          string    `gorm:"column:country" json:"country"`
	Gender           int       `gorm:"column:gender" json:"gender"`
	Language         string    `gorm:"column:language" json:"language"`
	NickName         string    `gorm:"column:nick_name" json:"nickName"`
	Province         string    `gorm:"column:province" json:"province"`
	RemainUsageCount int64     `gorm:"column:remain_usage_count" json:"remainUsageCount"`
	TotalUsageCount  int64     `gorm:"column:total_usage_count" json:"totalUsageCount"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
