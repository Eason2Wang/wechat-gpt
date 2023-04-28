package model

import (
	"time"

	"github.com/google/uuid"
)

// UserModel 用户模型
type OrderModel struct {
	Id             uuid.UUID `gorm:"column:id" json:"id"`
	UserId         string    `gorm:"column:user_id" json:"userId"`
	OpenId         string    `gorm:"column:open_id" json:"openId"`
	OutTradeNo     string    `gorm:"column:out_trade_no" json:"outTradeNo"`
	TotalFee       uint32    `gorm:"column:total_fee" json:"totalFee"`
	SpbillCreateIp string    `gorm:"column:spbill_create_ip" json:"spbillCreateIp"`
	Status         int8      `gorm:"column:status" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
