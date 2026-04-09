package service

import (
	"context"
	"errors"
	"finalai/internal/dto"
	"finalai/internal/repository"
	mycrypt "finalai/pkg/crypt"
	myjwt "finalai/pkg/jwt"
	"strings"
)

type UserSVC struct {
	UserDAO *repository.UserDAO
}

func NewUserSVC() *UserSVC {
	return &UserSVC{
		UserDAO: repository.NewUserDAO(),
	}
}

func (s *UserSVC) Register(ctx context.Context, req *dto.UserRegisterReq) (*dto.UserRegisterRes, error) {
	username := strings.TrimSpace(req.Username)
	password, err := mycrypt.GetHashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 查看用户是否已存在
	existingUser, err := s.UserDAO.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("failed to check existing user")
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	err = s.UserDAO.CreateUser(username, password)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return &dto.UserRegisterRes{
		Username: username,
	}, nil
}

func (s *UserSVC) Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserLoginRes, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	user, err := s.UserDAO.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("failed to query user")
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if err = mycrypt.CheckHashAndPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := myjwt.GenerateToken(myjwt.MyData{Username: user.Username})
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.UserLoginRes{
		Username: user.Username,
		Token:    token,
	}, nil
}
