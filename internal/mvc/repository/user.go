package repository

import (
	"context"
	"errors"

	"finalai/internal/database/mysql"
	"finalai/internal/model"

	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, user *model.User) error {
	return mysql.DB.WithContext(ctx).Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := mysql.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
