package model

import (
	"time"

	"github.com/google/uuid"
)

// UserModel 用户模型
type PromptModel struct {
	Id        uuid.UUID `gorm:"column:id" json:"id"`
	UserId    string    `gorm:"column:user_id" json:"userId"`
	Prompt    string    `gorm:"column:prompt" json:"prompt"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
