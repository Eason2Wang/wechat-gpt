package dao

import (
	"ininpop-chatgpt/db"
	"ininpop-chatgpt/db/model"
	"time"
)

const userTableName = "users"
const counterTableName = "Counters"
const reportTableName = "wechat_report"

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