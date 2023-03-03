package dao

import (
	"ininpop-chatgpt/db"
	"ininpop-chatgpt/db/model"
	"time"
)

const userTableName = "Users"
const counterTableName = "Counters"
const reportTableName = "wechat_report"

// ClearCounter 清除Counter
func (imp *CounterInterfaceImp) ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(counterTableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func (imp *CounterInterfaceImp) UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(counterTableName).Save(counter).Error
}

// GetCounter 查询Counter
func (imp *CounterInterfaceImp) GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(counterTableName).Where("id = ?", id).First(counter).Error

	return counter, err
}

// GetUserList 查询UserList
func (imp *UserInterfaceImp) GetUserList() ([]model.UserModel, error) {
	var err error
	var users []model.UserModel

	cli := db.Get()
	err = cli.Table(userTableName).Find(&users).Error
	return users, err
}

// UpsertUser 更新/写入user
func (imp *UserInterfaceImp) UpsertUser(user *model.UserModel) error {
	cli := db.Get()
	return cli.Table(userTableName).Save(user).Error
}

func (imp *UserInterfaceImp) GetUserByOpenId(openId string) (*model.UserModel, error) {
	var err error
	var user = new(model.UserModel)

	cli := db.Get()
	err = cli.Table(userTableName).Where("open_id = ?", openId).First(user).Error
	return user, err
}

// UpdateSubscribe 更新subscribe
func (imp *UserInterfaceImp) UpdateSubscribe(openId string, subscribe uint) error {
	cli := db.Get()
	return cli.Table(userTableName).Where("open_id = ?", openId).Updates(map[string]interface{}{"subscribe": subscribe, "updated_at": time.Now()}).Error
}

// UpdateFollow 更新follow
func (imp *UserInterfaceImp) UpdateFollow(openId string, follow uint) error {
	cli := db.Get()
	return cli.Table(userTableName).Where("open_id = ?", openId).Updates(map[string]interface{}{"follow": follow, "updated_at": time.Now()}).Error
}

func (imp *UserInterfaceImp) UpdateFollowAndSubscribe(openId string, follow uint, subscribe uint) error {
	cli := db.Get()
	return cli.Table(userTableName).Where("open_id = ?", openId).Updates(map[string]interface{}{"follow": follow, "subscribe": subscribe, "updated_at": time.Now()}).Error
}

// GetAllReports 查询所有报告
func (imp *DailyReportInterfaceImp) GetAllReports() ([]model.ReportModel, error) {
	var err error
	var reports []model.ReportModel

	cli := db.Get()
	err = cli.Table(reportTableName).Find(&reports).Error
	return reports, err
}

// GetLatestReport 查询最新的报告
func (imp *DailyReportInterfaceImp) GetLatestReport() (*model.ReportModel, error) {
	var err error
	var report = new(model.ReportModel)

	cli := db.Get()
	err = cli.Table(reportTableName).Order("report_date desc").First(report).Error
	return report, err
}

// GetLatestReport 查询最新的报告
func (imp *DailyReportInterfaceImp) GetReportByDate(date string, sort string) (*model.ReportModel, error) {
	var err error
	var report = new(model.ReportModel)

	cli := db.Get()
	err = cli.Table(reportTableName).Where("report_date = ?", date).Order("created_at " + sort).First(report).Error
	return report, err
}
