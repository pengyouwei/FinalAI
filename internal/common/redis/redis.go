package redis

import (
	"context"
	"finalai/internal/config"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var DB *redis.Client
var Rdb *redis.Client

func Init() {
	config := config.GetConfig().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:          config.Host + ":" + strconv.Itoa(config.Port),
		Password:      config.Password,
		DB:            config.DB,
		UnstableResp3: true,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to [Redis]: " + err.Error())
	}

	DB = rdb
	Rdb = rdb

	slog.Info("Successfully connected to [Redis]")
}

func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			slog.Error("Failed to close [Redis] connection: " + err.Error())
			return
		}
	}
	slog.Info("Successfully closed [Redis] connection")
}

// InitRedisIndex 初始化 Redis 索引，支持按文件名区分
func InitRedisIndex(ctx context.Context, filename string, dimension int) error {
	indexName := GenerateIndexName(filename)

	// 检查索引是否存在
	_, err := DB.Do(ctx, "FT.INFO", indexName).Result()
	if err == nil {
		fmt.Println("索引已存在，跳过创建")
		return nil
	}

	// 如果索引不存在，创建新索引
	if !isIndexNotFoundError(err) {
		return fmt.Errorf("检查索引失败: %w", err)
	}

	fmt.Println("正在创建 Redis 索引...")

	prefix := GenerateIndexNamePrefix(filename)

	// 创建索引
	createArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", prefix,
		"SCHEMA",
		"content", "TEXT",
		"metadata", "TEXT",
		"vector", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	if err := DB.Do(ctx, createArgs...).Err(); err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}

	fmt.Println("索引创建成功！")
	return nil
}

// DeleteRedisIndex 删除 Redis 索引，支持按文件名区分
func DeleteRedisIndex(ctx context.Context, filename string) error {
	indexName := GenerateIndexName(filename)

	// 删除索引
	if err := DB.Do(ctx, "FT.DROPINDEX", indexName).Err(); err != nil {
		if isIndexNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("删除索引失败: %w", err)
	}

	fmt.Println("索引删除成功！")
	return nil
}

func GenerateIndexName(filename string) string {
	indexName := fmt.Sprintf(config.DefaultRedisKeyConfig.IndexName, filename)
	return indexName
}

func GenerateIndexNamePrefix(filename string) string {
	prefix := fmt.Sprintf(config.DefaultRedisKeyConfig.IndexNamePrefix, filename)
	return prefix
}

func isIndexNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unknown index name") || strings.Contains(msg, "no such index")
}
