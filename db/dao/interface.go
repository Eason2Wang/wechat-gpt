package dao

import (
	"ininpop-chatgpt/db/model"
)

// UserInterface 用户数据模型接口
type UserInterface interface {
	GetUserList() ([]model.UserModel, error)
	GetUserByOpenId(openId string) (*model.UserModel, error)
	UpsertUser(user *model.UserModel) error
	UpdateSubscribe(openId string, subscribe uint) error
	UpdateFollow(openId string, follow uint) error
	UpdateFollowAndSubscribe(openId string, follow uint, subscribe uint) error
}

// UserInterfaceImp 用户数据模型实现
type UserInterfaceImp struct{}

// Imp 实现实例
var UserImp UserInterface = &UserInterfaceImp{}