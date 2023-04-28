package dao

import (
	"time"
	"wechat-gpt/db"
	"wechat-gpt/db/model"
)

const userTableName = "users"
const promptTableName = "prompts"
const orderTableName = "orders"

// GetUserList 查询UserList
func (imp *UserInterfaceImp) GetUserList() ([]model.UserModel, error) {
	var err error
	var users []model.UserModel

	cli := db.Get()
	err = cli.Table(userTableName).Find(&users).Error
	return users, err
}

// InsertUser 更新/写入user
func (imp *UserInterfaceImp) InsertUser(user *model.UserModel) error {
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

func (imp *UserInterfaceImp) UpdateNickNameAndAvatar(openId string, nickName string, avatar string) error {
	cli := db.Get()
	return cli.Table(userTableName).Where("open_id = ?", openId).Updates(map[string]interface{}{"nick_name": nickName, "avatar_url": avatar, "updated_at": time.Now()}).Error
}

func (imp *UserInterfaceImp) UpdateRemainUsage(openId string, usage int64) error {
	cli := db.Get()
	return cli.Table(userTableName).Where("open_id = ?", openId).Updates(map[string]interface{}{"remain_usage_count": usage, "updated_at": time.Now()}).Error
}

func (imp *UserInterfaceImp) UpdateFollowAndSubscribe(openId string, follow uint, subscribe uint) error {
	cli := db.Get()
	return cli.Table(userTableName).Where("open_id = ?", openId).Updates(map[string]interface{}{"follow": follow, "subscribe": subscribe, "updated_at": time.Now()}).Error
}

func (imp *PromptInterfaceImp) InsertPrompt(prompt *model.PromptModel) error {
	cli := db.Get()
	return cli.Table(promptTableName).Save(prompt).Error
}

func (imp *OrderInterfaceImp) InsertOrder(order *model.OrderModel) error {
	cli := db.Get()
	return cli.Table(orderTableName).Save(order).Error
}

func (imp *OrderInterfaceImp) GetOrderByTradeNo(TradeNo string) (*model.OrderModel, error) {
	var err error
	var order = new(model.OrderModel)

	cli := db.Get()
	err = cli.Table(orderTableName).Where("out_trade_no = ?", TradeNo).First(order).Error
	return order, err
}

func (imp *OrderInterfaceImp) UpdateOrderStatus(TradeNo string, Status int8) error {
	cli := db.Get()
	return cli.Table(orderTableName).Where("out_trade_no = ?", TradeNo).Updates(map[string]interface{}{"status": Status, "updated_at": time.Now()}).Error
}
