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
	MaxIdle  int    `toml:"maxIdle"`
	MaxOpen  int    `toml:"maxOpen"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type RabbitMQConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Vhost    string `toml:"vhost"`
}

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Mysql    MysqlConfig    `toml:"mysql"`
	Redis    RedisConfig    `toml:"redis"`
	RabbitMQ RabbitMQConfig `toml:"rabbitmq"`
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
