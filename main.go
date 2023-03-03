package main

import (
	"bytes"
	"context"
	"ininpop-chatgpt/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initGinServer()
}

func initGinServer() {
	var router *gin.Engine
	gin.SetMode(gin.DebugMode)
	router = gin.Default()

	service.ChatGptHandler(router)
	service.ChatGptStreamedHandler(router)
	service.UserHandler(router)

	srv := &http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server running...")
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// middleware app debug log, 作用是记录请求响应的信息
func AppDebugLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// record response info
		blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		ctx.Next()

		log.Println("responseInfo: ", blw.Body.String())
	}
}

// 初始化中间件
func InitMiddleware(r *gin.Engine) {
	//access log
	r.Use(AppDebugLog())
	// 异常保护
	r.Use(gin.Recovery())
}
