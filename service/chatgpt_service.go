package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ininpop-chatgpt/service/entity"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatGptHandler(router *gin.Engine) {
	router.POST("/chat", func(c *gin.Context) {
		chat(c)
	})
}

func chat(c *gin.Context) {
	var req entity.ChatGptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求数据错误: ", err)
		c.String(http.StatusBadRequest, "请求数据错误")
		return
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
		"http://43.153.59.188:5001/chat",
		"application/json",
		bytes.NewReader(bytesData),
	)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	// 	fmt.Printf("请求成功: %s", resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	data := make(map[string]interface{})
	json.Unmarshal(body, &data)
	fmt.Println("请求成功: ", data)
	bytes, _ := json.Marshal(data)
	c.JSON(http.StatusOK, string(bytes))
}