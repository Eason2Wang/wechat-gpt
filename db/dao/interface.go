package dao

import (
	"wechat-gpt/db/model"
)

// UserInterface 用户数据模型接口
type UserInterface interface {
	GetUserList() ([]model.UserModel, error)
	GetUserByOpenId(openId string) (*model.UserModel, error)
	InsertUser(user *model.UserModel) error
	UpdateNickNameAndAvatar(openId string, nickName string, avatar string) error
	UpdateRemainUsage(openId string, usage int64) error
	UpdateTotalUsage(openId string, usage int64) error
	UpdateFollowAndSubscribe(openId string, follow uint, subscribe uint) error
}

type PromptInterface interface {
	InsertPrompt(prompt *model.PromptModel) error
}

type OrderInterface interface {
	GetOrderByTradeNo(TradeNo string) (*model.OrderModel, error)
	InsertOrder(order *model.OrderModel) error
	UpdateOrderStatus(TradeNo string, Status int8) error
}

// UserInterfaceImp 用户数据模型实现
type UserInterfaceImp struct{}
type PromptInterfaceImp struct{}
type OrderInterfaceImp struct{}

// Imp 实现实例
var UserImp UserInterface = &UserInterfaceImp{}
var PromptImp PromptInterface = &PromptInterfaceImp{}
var OrderImp OrderInterface = &OrderInterfaceImp{}