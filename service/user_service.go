package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"wechat-gpt/db/dao"
	"wechat-gpt/db/model"
	"wechat-gpt/service/entity"
	"wechat-gpt/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserHandler 用户接口
func UserHandler(router *gin.Engine) {

	router.POST("/api/updateInfo", func(c *gin.Context) {
		httpCode, result := updateInfo(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/chat", func(c *gin.Context) {
		httpCode, result := chat(c)
		c.JSON(httpCode, result)
	})
}

func chat(c *gin.Context) (int, entity.Response) {
	var req entity.ChatGptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求数据错误: ", err)
		c.String(http.StatusBadRequest, "请求数据错误")
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	fmt.Println("请求数据: ", req)

	params := make(map[string]interface{})
	params["prompt"] = req.Prompt
	if req.ParentId != "" {
		params["parent_id"] = req.ParentId
	}
	if req.ConversationId != "" {
		params["conversation_id"] = req.ConversationId
	}
	fmt.Println("发送数据: ", params)
	bytesData, _ := json.Marshal(params)
	fmt.Println("发送数据json: ", bytesData)
	resp, err := http.Post(
		"https://chatgpt-api.ininpop.com/chatgpt-api/chat",
		"application/json",
		bytes.NewReader(bytesData),
	)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return http.StatusInternalServerError, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	// 	fmt.Printf("请求成功: %s", resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return http.StatusInternalServerError, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	data := make(map[string]interface{})
	json.Unmarshal(body, &data)
	fmt.Println("请求成功: ", string(body))
	return http.StatusOK, entity.Response{
		Code: 0,
		Data: string(body),
	}
}

// login 获取并保存用户信息
func updateInfo(c *gin.Context) (int, entity.Response) {
	var req entity.UpdateInfoRequest
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
	userInfo, err := upsertUser(&req, openid)
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
func upsertUser(req *entity.UpdateInfoRequest, openid string) (*model.UserModel, error) {
	currentUser, err := dao.UserImp.GetUserByOpenId(openid)
	var user model.UserModel
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		user = model.UserModel{
			Id:         uuid.New(),
			OpenId:     openid,
			AvatarUrl:  req.AvatarUrl,
			City:       req.City,
			Country:    req.Country,
			Gender:     req.Gender,
			Language:   req.Language,
			NickName:   req.NickName,
			Province:   req.Province,
			UsageCount: 10,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	} else {
		user = model.UserModel{
			Id:         currentUser.Id,
			OpenId:     currentUser.OpenId,
			AvatarUrl:  req.AvatarUrl,
			City:       req.City,
			Country:    req.Country,
			Gender:     req.Gender,
			Language:   req.Language,
			NickName:   req.NickName,
			Province:   req.Province,
			UsageCount: req.UsageCount,
			CreatedAt:  currentUser.CreatedAt,
			UpdatedAt:  time.Now(),
		}
	}
	err = dao.UserImp.UpsertUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
