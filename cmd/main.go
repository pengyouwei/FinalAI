package main

import (
	"finalai/internal/config"
	db "finalai/internal/database"
	myjwt "finalai/pkg/jwt"
	mylogger "finalai/pkg/logger"
	myvalidator "finalai/pkg/validator"

	"finalai/internal/model"
	"finalai/internal/mq/rabbitmq"
	"finalai/internal/router"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/labstack/echo/v5"
)

func init() {
	config.Init()
}

func main() {
	// 初始化数据库和消息队列连接
	db.InitMysql()
	db.InitRedis()
	rabbitmq.Init()
	myjwt.Init()

	// 迁移表
	db.MysqlDB.AutoMigrate(&model.User{}, &model.Session{}, &model.Message{})

	// 启动HTTP服务器
	e := echo.New()
	e.Logger = mylogger.NewLogger()
	e.Validator = myvalidator.NewValidator()
	router.RegisterRoutes(e)

	go func() {
		if err := e.Start(":" + strconv.Itoa(config.GetConfig().Server.Port)); err != nil {
			e.Logger.Warn("Failed to start server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Shutting down server...")
}
