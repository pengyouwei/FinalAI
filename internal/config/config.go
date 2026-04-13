package config

import "github.com/BurntSushi/toml"

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type MysqlConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	LogLevel string `toml:"log_level"`
	Charset  string `toml:"charset"`
	MaxIdle  int    `toml:"max_idle"`
	MaxOpen  int    `toml:"max_open"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type RabbitMQConfig struct {
	Host          string `toml:"host"`
	Port          int    `toml:"port"`
	Username      string `toml:"username"`
	Password      string `toml:"password"`
	Vhost         string `toml:"vhost"`
	ConsumerCount int    `toml:"consumer_count"`
}

type RagModelConfig struct {
	RagBaseUrl        string `toml:"rag_base_url"`
	RagChatModelName  string `toml:"rag_chat_model_name"`
	RagEmbeddingModel string `toml:"rag_embedding_model"`
	RagDimension      int    `toml:"rag_dimension"`
}

type Config struct {
	Server         ServerConfig   `toml:"server"`
	Mysql          MysqlConfig    `toml:"mysql"`
	Redis          RedisConfig    `toml:"redis"`
	RabbitMQ       RabbitMQConfig `toml:"rabbitmq"`
	RagModelConfig RagModelConfig `toml:"rag_model"`
}

type RedisKeyConfig struct {
	CaptchaPrefix   string
	IndexName       string
	IndexNamePrefix string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	CaptchaPrefix:   "captcha:%s",
	IndexName:       "rag_docs:%s:idx",
	IndexNamePrefix: "rag_docs:%s:",
}

var config *Config

func Init() {
	config = &Config{}
	if _, err := toml.DecodeFile("./config/config.toml", config); err != nil {
		panic("Failed to load config: " + err.Error())
	}
}

func GetConfig() *Config {
	if config == nil {
		Init()
	}
	return config
}
