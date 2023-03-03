package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ininpop-chatgpt/db/dao"
	"ininpop-chatgpt/db/model"
	"ininpop-chatgpt/service/entity"
	"ininpop-chatgpt/utils"
	"io/ioutil"
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

// GetRequestBodyJson 获取JSON请求体参数
func GetRequestBodyJson(r *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

func postGetOpenData(appid string, openid string, cloudid string) map[string]interface{} {
	params := make(map[string]interface{})
	params["cloudid_list"] = [1]string{cloudid}
	bytesData, _ := json.Marshal(params)
	url := fmt.Sprintf("http://api.weixin.qq.com/wxa/getopendata?from_appid=%s&openid=%s", appid, openid)
	fmt.Println("OpenDataUrl:", url)
	resp, _ := http.Post(
		url,
		"application/json",
		bytes.NewReader(bytesData),
	)
	body, _ := ioutil.ReadAll(resp.Body)
	bodyMap, err := JsonToMap(string(body))
	if err == nil {
		fmt.Println("OpenDataBodyMap:", bodyMap)
		return bodyMap
	}
	return nil
}

// Convert json string to map
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}

	return m, nil
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
