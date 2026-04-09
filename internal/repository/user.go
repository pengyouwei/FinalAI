package repository

import (
	"errors"
	db "finalai/internal/database"
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
	return db.MysqlDB.Create(user).Error
}

func (dao *UserDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := db.MysqlDB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
