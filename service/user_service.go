package service

import (
	"fmt"
	"wechat-gpt/db/dao"
	"wechat-gpt/db/model"
	"wechat-gpt/service/entity"
	"wechat-gpt/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserHandler 用户接口
func UserHandler(router *gin.Engine) {
	router.POST("/login", func(c *gin.Context) {
		httpCode, result := login(c)
		c.JSON(httpCode, result)
	})
}

// login 获取并保存用户信息
func login(c *gin.Context) (int, entity.Response) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}

	var userInfo *model.UserModel

	header := c.Request.Header
	fmt.Println("Header全部数据:", header)
	var appid string
	if header["X-Wx-From-Appid"] == nil {
		appid = ""
	} else {
		appid = header["X-Wx-From-Appid"][0]
	}
	fmt.Println("appid:", appid)
	var openid string
	if header["X-Wx-From-Openid"] != nil {
		openid = header["X-Wx-From-Openid"][0]
	} else if header["X-Wx-Openid"] != nil {
		openid = header["X-Wx-Openid"][0]
	} else {
		//只有在测试环境下才有这个参数
		openid = req.OpenId
	}
	fmt.Println("openid:", openid)
	userInfo, err := upsertUser(&entity.UserInfo{
		OpenId: openid,
	})
	if err != nil {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	} else {
		return http.StatusOK, entity.Response{
			Code: 0,
			Data: userInfo,
		}
	}
}

// upsertUser 更新或修改用户信息
func upsertUser(userInfo *entity.UserInfo) (*model.UserModel, error) {
	currentUser, err := dao.UserImp.GetUserByOpenId(userInfo.OpenId)
	var user model.UserModel
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		user = model.UserModel{
			Id:        uuid.New(),
			OpenId:    userInfo.OpenId,
			AvatarUrl: userInfo.AvatarUrl,
			City:      userInfo.City,
			Country:   userInfo.Country,
			Gender:    userInfo.Gender,
			Language:  userInfo.Language,
			NickName:  userInfo.NickName,
			Province:  userInfo.Province,
			AppId:     userInfo.WaterMark.AppId,
			Follow:    1,
			Subscribe: 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	} else {
		user = model.UserModel{
			Id:        currentUser.Id,
			OpenId:    currentUser.OpenId,
			AvatarUrl: userInfo.AvatarUrl,
			City:      userInfo.City,
			Country:   userInfo.Country,
			Gender:    userInfo.Gender,
			Language:  userInfo.Language,
			NickName:  userInfo.NickName,
			Province:  userInfo.Province,
			AppId:     userInfo.WaterMark.AppId,
			Follow:    currentUser.Follow,
			Subscribe: currentUser.Subscribe,
			CreatedAt: currentUser.CreatedAt,
			UpdatedAt: time.Now(),
		}
	}
	err = dao.UserImp.UpsertUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
