package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"wechat-gpt/db/dao"
	"wechat-gpt/db/model"
	"wechat-gpt/service/entity"
	"wechat-gpt/utils"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/green"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserHandler 用户接口
func UserHandler(router *gin.Engine) {

	router.POST("/api/checkExist", func(c *gin.Context) {
		httpCode, result := checkExist(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/generateUser", func(c *gin.Context) {
		httpCode, result := generateUser(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/updateInfo", func(c *gin.Context) {
		httpCode, result := updateNickNameAndAvatar(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/updateUsage", func(c *gin.Context) {
		httpCode, result := updateRemainUsage(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/getUserInfo", func(c *gin.Context) {
		httpCode, result := getUserInfo(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/checkText", func(c *gin.Context) {
		httpCode, result := checkText(c)
		c.JSON(httpCode, result)
	})
}

func getOpenId(c *gin.Context) string {
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
	}
	fmt.Println("openid:", openid)
	return openid
}

// checkExist 检查用户是否存在
func checkExist(c *gin.Context) (int, entity.Response) {
	var req entity.EmptyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	openid := getOpenId(c)
	currentUser, err := dao.UserImp.GetUserByOpenId(openid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	} else if err == gorm.ErrRecordNotFound {
		return http.StatusOK, entity.Response{
			Code:     utils.USER_NOT_FOUNT,
			ErrorMsg: "user not found",
		}
	} else {
		currentUser.OpenId = ""
		return http.StatusOK, entity.Response{
			Code: utils.SUCCESS,
			Data: currentUser,
		}
	}
}

// generateUser 生成新用户
func generateUser(c *gin.Context) (int, entity.Response) {
	var req entity.GenerateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	openid := getOpenId(c)
	_, err := dao.UserImp.GetUserByOpenId(openid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	} else if err == gorm.ErrRecordNotFound {
		user := model.UserModel{
			Id:               uuid.New(),
			OpenId:           openid,
			AvatarUrl:        req.AvatarUrl,
			City:             req.City,
			Country:          req.Country,
			Gender:           req.Gender,
			Language:         req.Language,
			NickName:         req.NickName,
			Province:         req.Province,
			RemainUsageCount: 10,
			TotalUsageCount:  10,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		err = dao.UserImp.InsertUser(&user)
		if err != nil {
			return http.StatusOK, entity.Response{
				Code:     utils.SERVER_DB_ERR,
				ErrorMsg: err.Error(),
			}
		}
		return http.StatusOK, entity.Response{
			Code: 0,
			Data: user,
		}
	} else {
		return http.StatusOK, entity.Response{
			Code:     utils.USER_ALREADY_EXIST,
			ErrorMsg: "user already exist",
		}
	}
}

// updateNickNameAndAvatar
func updateNickNameAndAvatar(c *gin.Context) (int, entity.Response) {
	var req entity.UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	openid := getOpenId(c)
	currentUser, err := dao.UserImp.GetUserByOpenId(openid)
	if err != nil {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	} else {
		var avatarUrl string
		if req.AvatarUrl == "" {
			avatarUrl = currentUser.AvatarUrl
		} else {
			avatarUrl = req.AvatarUrl
		}
		err = dao.UserImp.UpdateNickNameAndAvatar(openid, req.NickName, avatarUrl)
		if err != nil {
			return http.StatusOK, entity.Response{
				Code:     utils.SERVER_DB_ERR,
				ErrorMsg: err.Error(),
			}
		}
		return http.StatusOK, entity.Response{
			Code: 0,
			Data: nil,
		}
	}
}

func updateRemainUsage(c *gin.Context) (int, entity.Response) {
	var req entity.UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	openid := getOpenId(c)
	user, err := dao.UserImp.GetUserByOpenId(openid)
	if err != nil {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	} else {
		err = dao.UserImp.UpdateRemainUsage(openid, req.RemainUsageCount)
		if err != nil {
			return http.StatusOK, entity.Response{
				Code:     utils.SERVER_DB_ERR,
				ErrorMsg: err.Error(),
			}
		}
		user.OpenId = ""
		user.RemainUsageCount = req.RemainUsageCount
		return http.StatusOK, entity.Response{
			Code: 0,
			Data: user,
		}
	}
}

func getUserInfo(c *gin.Context) (int, entity.Response) {
	var req entity.EmptyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	openid := getOpenId(c)
	user, err := dao.UserImp.GetUserByOpenId(openid)
	if err != nil {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	} else {
		user.OpenId = ""
		return http.StatusOK, entity.Response{
			Code: 0,
			Data: user,
		}
	}
}

func checkText(c *gin.Context) (int, entity.Response) {
	var req entity.CheckTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	client, _err := green.NewClientWithAccessKey("cn-shenzhen", os.Getenv("ALIYUN_ACCESS_KEY"), os.Getenv("ALIYUN_ACCESS_SECRET"))
	if _err != nil {
		fmt.Println(_err.Error())
		return http.StatusOK, entity.Response{
			Code:     utils.ALIYUN_GREEN_ERR,
			ErrorMsg: _err.Error(),
		}
	}
	task := map[string]interface{}{"content": req.Text}
	// scenes：检测场景，唯一取值：antispam。
	content, _ := json.Marshal(
		map[string]interface{}{
			"scenes": [...]string{"antispam"},
			"tasks":  [...]map[string]interface{}{task},
		},
	)

	textScanRequest := green.CreateTextScanRequest()
	textScanRequest.SetContent(content)
	textScanResponse, err := client.TextScan(textScanRequest)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusOK, entity.Response{
			Code:     utils.ALIYUN_GREEN_ERR,
			ErrorMsg: _err.Error(),
		}
	}
	if textScanResponse.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(textScanResponse.GetHttpStatus()))
		return http.StatusOK, entity.Response{
			Code:     utils.ALIYUN_GREEN_ERR,
			ErrorMsg: "response not success. status:" + strconv.Itoa(textScanResponse.GetHttpStatus()),
		}
	}
	data := make(map[string]interface{})
	json.Unmarshal(textScanResponse.GetHttpContentBytes(), &data)
	fmt.Println(data)
	return http.StatusOK, entity.Response{
		Code: 0,
		Data: data,
	}
}
