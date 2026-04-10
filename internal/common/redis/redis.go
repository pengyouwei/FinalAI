package redis

import (
	"finalai/internal/config"
	"log/slog"
	"strconv"

	"github.com/go-redis/redis"
)

var DB *redis.Client

func Init() {
	config := config.GetConfig().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		panic("Failed to connect to [Redis]: " + err.Error())
	}

	DB = rdb

	slog.Info("Successfully connected to [Redis]")
}
