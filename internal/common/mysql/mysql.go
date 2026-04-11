package mysql

import (
	"finalai/internal/config"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	config := config.GetConfig().Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect to [Mysql]: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get SQL DB instance: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	slog.Info("Successfully connected to [Mysql]")
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		slog.Error("Failed to get SQL DB instance: " + err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		slog.Error("Failed to close SQL DB connection: " + err.Error())
	}
	slog.Info("Successfully closed [Mysql] connection")
}
