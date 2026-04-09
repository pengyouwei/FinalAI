package db

import (
	"finalai/internal/config"
	"log/slog"
	"strconv"

	"github.com/go-redis/redis"
)

var RedisDB *redis.Client

func InitRedis() {
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

	RedisDB = rdb

	slog.Info("Successfully connected to [Redis]")
}
