package repository

import (
	"errors"

	"finalai/internal/database/mysql"
	"finalai/internal/model"

	"gorm.io/gorm"
)

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (dao *UserDAO) CreateUser(username, password string) error {
	user := &model.User{
		Username: username,
		Password: password,
	}
	return mysql.DB.Create(user).Error
}

func (dao *UserDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := mysql.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
