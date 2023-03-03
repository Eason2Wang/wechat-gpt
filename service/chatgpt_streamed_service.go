package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"ininpop-chatgpt/service/entity"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ChatGptStreamedHandler(router *gin.Engine) {
	router.POST("/stream", func(c *gin.Context) {
		stream(c)
	})

	// stream := NewServer()
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 10)
	// 		now := time.Now().Format("2006-01-02 15:04:05")
	// 		currentTime := fmt.Sprintf("The Current Time Is %v", now)
	// 		fmt.Println("currentTime: ", currentTime)
	// 		// Send current time to clients message channel
	// 		stream.Message <- currentTime
	// 	}
	// }()
	// router.GET("/stream_test", HeadersMiddleware(), stream.serveHTTP(), func(c *gin.Context) {
	// 	v, ok := c.Get("clientChan")
	// 	fmt.Println("ok1: ", ok)
	// 	if !ok {
	// 		return
	// 	}
	// 	clientChan, ok := v.(ClientChan)
	// 	fmt.Println("ok2: ", ok)
	// 	if !ok {
	// 		return
	// 	}
	// 	c.Stream(func(w io.Writer) bool {
	// 		// Stream message to client from message channel
	// 		if msg, ok := <-clientChan; ok {
	// 			fmt.Println("msg: ", msg)
	// 			c.SSEvent("message", msg)
	// 			return true
	// 		}
	// 		return false
	// 	})
	// })
	router.POST("/stream_test", func(c *gin.Context) {
		stream_test(c)
	})
}

func stream(c *gin.Context) {
	var req entity.ChatGptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求数据错误: ", err)
		c.String(http.StatusBadRequest, "请求数据错误")
		return
	}
	fmt.Println("请求数据: ", req)

	params := make(map[string]interface{})
	params["prompt"] = req.Prompt
	params["streamed"] = true
	if req.ParentId != "" {
		params["parent_id"] = req.ParentId
	}
	if req.ConversationId != "" {
		params["conversation_id"] = req.ConversationId
	}
	fmt.Println("发送数据: ", params)
	bytesData, _ := json.Marshal(params)
	fmt.Println("发送数据json: ", bytesData)
	request, err := http.NewRequest(http.MethodPost, "http://43.153.59.188:5001/chat", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	ch := make(chan string)

	bufferedReader := bufio.NewReader(resp.Body)

	// Read the response body
	go func() {
		for {
			buffer := make([]byte, 4*1024)
			len, err := bufferedReader.Read(buffer)
			if len > 0 {
				splitStr := strings.Split(string(buffer), "%$")
				for _, s := range splitStr {
					data := make(map[string]interface{})
					json.Unmarshal([]byte(s), &data)
					fmt.Println("Response: ", data)
					bytes, _ := json.Marshal(data)
					ch <- string(bytes)
				}
			}
			buffer = nil
			if err != nil {
				if err == io.EOF {
				}
				close(ch)
				break
			}
		}
	}()

	clientGone := c.Request.Context().Done()
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-ch:
			if ok {
				c.SSEvent("message", msg)
			}
			return ok
		case <-clientGone:
			fmt.Println("stream end")
			return false
		}
	})
}

// It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan string

func stream_test(c *gin.Context) {
	var req entity.ChatGptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("请求数据错误: ", err)
		c.String(http.StatusBadRequest, "请求数据错误")
		return
	}
	fmt.Println("请求数据: ", req)

	params := make(map[string]interface{})
	params["prompt"] = req.Prompt
	params["streamed"] = true
	if req.ParentId != "" {
		params["parent_id"] = req.ParentId
	}
	if req.ConversationId != "" {
		params["conversation_id"] = req.ConversationId
	}
	fmt.Println("发送数据: ", params)
	bytesData, _ := json.Marshal(params)
	fmt.Println("发送数据json: ", bytesData)
	request, err := http.NewRequest(http.MethodPost, "http://43.153.59.188:5001/stream_test", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
	}
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("请求失败: %s", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	ch := make(chan string)

	bufferedReader := bufio.NewReader(resp.Body)

	var totalBytesReceived int

	// Read the response body
	go func() {
		for {
			buffer := make([]byte, 4*1024)
			len, err := bufferedReader.Read(buffer)
			if len > 0 {
				totalBytesReceived += len
				log.Println(len, "bytes received")
				fmt.Println("Response: Read", string(buffer))
				ch <- string(buffer)
			}
			buffer = nil
			if err != nil {
				if err == io.EOF {
				}
				close(ch)
				break
			}
		}
	}()

	c.Stream(func(w io.Writer) bool {
		msg, ok := <-ch
		if ok {
			fmt.Println("ch: ", msg)
			c.SSEvent("message", msg)
			// c.String(http.StatusOK, msg)
		}
		return ok
	})
}

// Initialize event and Start procnteessing requests
func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
