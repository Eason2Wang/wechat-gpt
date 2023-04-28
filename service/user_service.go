package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"wechat-gpt/db/dao"
	"wechat-gpt/db/model"
	"wechat-gpt/service/entity"
	"wechat-gpt/utils"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alimt"
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

	router.POST("/api/translate", func(c *gin.Context) {
		httpCode, result := translate(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/savePrompt", func(c *gin.Context) {
		httpCode, result := savePrompt(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/unifiedOrder", func(c *gin.Context) {
		httpCode, result := unifiedOrder(c)
		c.JSON(httpCode, result)
	})

	router.POST("/api/payCallback", func(c *gin.Context) {
		httpCode, result := payCallback(c)
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
	client, _err := green.NewClientWithAccessKey(
		"cn-shenzhen",
		os.Getenv("ALIYUN_ACCESS_KEY"),
		os.Getenv("ALIYUN_ACCESS_SECRET"))
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

func translate(c *gin.Context) (int, entity.Response) {
	var req entity.TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	// 创建ecsClient实例
	alimtClient, err := alimt.NewClientWithAccessKey(
		"cn-hangzhou",                     // 地域ID
		os.Getenv("ALIYUN_ACCESS_KEY"),    // 您的Access Key ID
		os.Getenv("ALIYUN_ACCESS_SECRET")) // 您的Access Key Secret
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusOK, entity.Response{
			Code:     utils.ALIYUN_TRANSLATE_ERR,
			ErrorMsg: err.Error(),
		}
	}
	// 创建API请求并设置参数
	request := alimt.CreateTranslateECommerceRequest()
	// 等价于 request.PageSize = "10"
	request.Method = "POST"                 //设置请求
	request.FormatType = "text"             //翻译文本的格式
	request.SourceLanguage = req.SourceLang //源语言
	request.SourceText = req.Text           //原文
	request.TargetLanguage = req.TargetLang //目标语言
	request.Scene = "general"               // 通用版本默认是：general
	// 发起请求并处理异常
	response, err := alimtClient.TranslateECommerce(request)
	if err != nil {
		fmt.Println(err.Error())
		return http.StatusOK, entity.Response{
			Code:     utils.ALIYUN_TRANSLATE_ERR,
			ErrorMsg: err.Error(),
		}
	}
	fmt.Println(response)

	if response.GetHttpStatus() != 200 {
		fmt.Println("response not success. status:" + strconv.Itoa(response.GetHttpStatus()))
		return http.StatusOK, entity.Response{
			Code:     utils.ALIYUN_TRANSLATE_ERR,
			ErrorMsg: "response not success. status:" + strconv.Itoa(response.GetHttpStatus()),
		}
	}
	data := make(map[string]interface{})
	json.Unmarshal(response.GetHttpContentBytes(), &data)
	fmt.Println(data)
	return http.StatusOK, entity.Response{
		Code: 0,
		Data: data,
	}
}

func savePrompt(c *gin.Context) (int, entity.Response) {
	var req entity.SavePromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	prompt := model.PromptModel{
		Id:        uuid.New(),
		UserId:    req.UserId,
		Prompt:    req.Prompt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := dao.PromptImp.InsertPrompt(&prompt)
	if err != nil {
		return http.StatusOK, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	}
	return http.StatusOK, entity.Response{
		Code: 0,
		Data: prompt,
	}
}

func unifiedOrder(c *gin.Context) (int, entity.Response) {
	var req entity.UnifiedOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.Response{
			Code:     utils.SERVER_MISSING_PARAMS,
			ErrorMsg: err.Error(),
		}
	}
	header := c.Request.Header
	var openid string
	if header["X-Wx-From-Openid"] != nil {
		openid = header["X-Wx-From-Openid"][0]
	} else if header["X-Wx-Openid"] != nil {
		openid = header["X-Wx-Openid"][0]
	}
	params := make(map[string]interface{})
	params["callback_type"] = 2
	params["env_id"] = "prod-7g96cjuxd039a0e4"
	params["function_name"] = "payCallback"
	container := make(map[string]interface{})
	container["service"] = "payCallback"
	container["path"] = "/api/"
	params["container"] = container
	params["sub_mch_id"] = "1642970169"
	params["nonce_str"] = utils.RandomString(32)
	params["body"] = req.Body
	outTradeNo := utils.GenerateOrderNoByTime()
	params["out_trade_no"] = outTradeNo
	params["total_fee"] = req.TotalFee
	ip := header["X-Original-Forwarded-For"][0]
	params["spbill_create_ip"] = ip
	params["trade_type"] = "JSAPI"
	params["openid"] = openid

	// 先插入数据库
	order := model.OrderModel{
		Id:             uuid.New(),
		UserId:         req.UserId,
		OpenId:         openid,
		OutTradeNo:     outTradeNo,
		TotalFee:       req.TotalFee,
		SpbillCreateIp: ip,
		Status:         -1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := dao.OrderImp.InsertOrder(&order)
	if err != nil {
		return http.StatusInternalServerError, entity.Response{
			Code:     utils.SERVER_DB_ERR,
			ErrorMsg: err.Error(),
		}
	}

	bytesData, _ := json.Marshal(params)
	resp, err := http.Post(
		"http://api.weixin.qq.com/_/pay/unifiedorder",
		"application/json",
		bytes.NewReader(bytesData),
	)
	if err != nil {
		return http.StatusInternalServerError, entity.Response{
			Code:     utils.Failure,
			ErrorMsg: err.Error(),
		}
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败: %s", resp.Body)
		return http.StatusInternalServerError, entity.Response{
			Code:     utils.Failure,
			ErrorMsg: err.Error(),
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		return http.StatusInternalServerError, entity.Response{
			Code:     utils.Failure,
			ErrorMsg: err.Error(),
		}
	}
	data := make(map[string]interface{})
	json.Unmarshal(body, &data)
	fmt.Println("请求成功: ", data)
	return http.StatusOK, entity.Response{
		Code: 0,
		Data: data,
	}
}

var mutex sync.RWMutex

func payCallback(c *gin.Context) (int, entity.PayCallbackResponse) {
	var req entity.PayCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return http.StatusBadRequest, entity.PayCallbackResponse{
			ErrCode: utils.SERVER_MISSING_PARAMS,
			ErrMsg:  err.Error(),
		}
	}
	mutex.RLock() // 加读锁
	order, err := dao.OrderImp.GetOrderByTradeNo(req.OutTradeNo)
	mutex.RUnlock() // 释放读锁
	if err != nil {
		fmt.Println("获取数据库订单失败: ", err.Error())
		return http.StatusOK, entity.PayCallbackResponse{
			ErrCode: utils.SERVER_DB_ERR,
			ErrMsg:  err.Error(),
		}
	}
	if req.ReturnCode == "SUCCESS" {
		if req.ResultCode == "SUCCESS" {
			if req.TotalFee == order.TotalFee {
				mutex.Lock()
				dao.OrderImp.UpdateOrderStatus(req.OutTradeNo, 0)
				mutex.Unlock()
				fmt.Println("交易成功！")
				return http.StatusOK, entity.PayCallbackResponse{
					ErrCode: utils.SUCCESS,
					ErrMsg:  req.ReturnMsg,
				}
			} else {
				mutex.Lock()
				dao.OrderImp.UpdateOrderStatus(req.OutTradeNo, 1)
				mutex.Unlock()
				fmt.Println("返回金额与数据库金额不一致")
				return http.StatusOK, entity.PayCallbackResponse{
					ErrCode: utils.Failure,
					ErrMsg:  "返回金额与数据库金额不一致",
				}
			}
		} else {
			mutex.Lock()
			dao.OrderImp.UpdateOrderStatus(req.OutTradeNo, 1)
			mutex.Unlock()
			fmt.Println("callback ResultCode失败:", req.ErrCodeDes)
			return http.StatusOK, entity.PayCallbackResponse{
				ErrCode: utils.Failure,
				ErrMsg:  req.ErrCodeDes,
			}
		}
	} else {
		fmt.Println("callback ReturnCode失败:", req.ReturnMsg)
		return http.StatusOK, entity.PayCallbackResponse{
			ErrCode: utils.Failure,
			ErrMsg:  req.ReturnMsg,
		}
	}
}
