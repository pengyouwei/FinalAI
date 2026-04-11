package main

import (
	"context"
	"finalai/internal/common/aihelper"
	"finalai/internal/common/mysql"
	"finalai/internal/common/rabbitmq"
	"finalai/internal/common/redis"
	"finalai/internal/config"
	"finalai/internal/model"
	"finalai/internal/repository"
	"finalai/internal/router"
	myjwt "finalai/pkg/jwt"
	mylogger "finalai/pkg/logger"
	myvalidator "finalai/pkg/validator"

	"log/slog"

	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/labstack/echo/v5"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	config.Init()
}

func main() {
	myjwt.Init()

	mysql.Init()
	defer mysql.Close()

	redis.Init()
	defer redis.Close()

	rabbitmq.Init()
	defer rabbitmq.CloseConn()

	consumer := rabbitmq.StartConsumer(
		rabbitmq.MessageQueueName,
		func(msg *amqp.Delivery) error {
			return rabbitmq.HandleMessage(context.Background(), msg)
		},
	)
	defer consumer.Close()

	// 迁移表
	mysql.DB.AutoMigrate(&model.User{}, &model.Session{}, &model.Message{})

	// 启动HTTP服务器
	e := echo.New()
	e.Logger = mylogger.NewLogger()
	e.Validator = myvalidator.NewValidator()
	router.RegisterRoutes(e)

	// 从数据库加载消息并初始化 AIHelperManager
	readDataFromDB()

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

// 从数据库加载消息并初始化 AIHelperManager
func readDataFromDB() error {
	manager := aihelper.GetGlobalManager()
	// 从数据库读取所有消息
	msgs, err := repository.GetAllMessages()
	if err != nil {
		return err
	}
	// 遍历数据库消息
	for i := range msgs {
		m := msgs[i]
		//默认openai模型
		modelType := "1"
		config := make(map[string]any)

		// 创建对应的 AIHelper
		helper, err := manager.GetOrCreateAIHelper(m.Username, m.SessionID, modelType, config)
		if err != nil {
			slog.Error("Failed to create AIHelper for user: " + m.Username + ", session: " + m.SessionID + ", error: " + err.Error())
			continue
		}
		slog.Info("readDataFromDB init: " + helper.SessionID)
		// 添加消息到内存中(不开启存储功能)
		helper.AddMessage(m.Content, m.Username, m.IsUser, false)
	}

	slog.Info("AIHelperManager init successfully with messages from DB")
	return nil
}
