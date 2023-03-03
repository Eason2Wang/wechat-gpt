package dao

import (
	"ininpop-chatgpt/db/model"
)

// CounterInterface 计数器数据模型接口
type CounterInterface interface {
	GetCounter(id int32) (*model.CounterModel, error)
	UpsertCounter(counter *model.CounterModel) error
	ClearCounter(id int32) error
}

// UserInterface 用户数据模型接口
type UserInterface interface {
	GetUserList() ([]model.UserModel, error)
	GetUserByOpenId(openId string) (*model.UserModel, error)
	UpsertUser(user *model.UserModel) error
	UpdateSubscribe(openId string, subscribe uint) error
	UpdateFollow(openId string, follow uint) error
	UpdateFollowAndSubscribe(openId string, follow uint, subscribe uint) error
}

// DailyReportInterface 每日报告模型接口
type DailyReportInterface interface {
	GetAllReports() ([]model.ReportModel, error)
	GetLatestReport() (*model.ReportModel, error)
	GetReportByDate(date string, sort string) (*model.ReportModel, error)
}

// CounterInterfaceImp 计数器数据模型实现
type CounterInterfaceImp struct{}

// UserInterfaceImp 用户数据模型实现
type UserInterfaceImp struct{}

// DailyReportInterfaceImp 每日报告数据模型实现
type DailyReportInterfaceImp struct{}

// Imp 实现实例
var CounterImp CounterInterface = &CounterInterfaceImp{}
var UserImp UserInterface = &UserInterfaceImp{}
var ReportImp DailyReportInterface = &DailyReportInterfaceImp{}
