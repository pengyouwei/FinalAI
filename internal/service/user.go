package service

import (
	"context"
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/model"
	"finalai/internal/repository"
	mycrypt "finalai/pkg/crypt"
	myjwt "finalai/pkg/jwt"
	"strings"
)

type UserSVC struct {
}

func NewUserSVC() *UserSVC {
	return &UserSVC{}
}

func (s *UserSVC) Register(ctx context.Context, req *dto.UserRegisterReq) (*dto.UserRegisterRes, error) {
	username := strings.TrimSpace(req.Username)
	password, err := mycrypt.GetHashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		return nil, apperror.ErrInternal.WithDetail("密码加密失败")
	}

	// 查看用户是否已存在
	existingUser, err := repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, apperror.ErrInternal.WithDetail("查询用户失败")
	}
	if existingUser != nil {
		return nil, apperror.ErrUserAlreadyExists
	}

	user := &model.User{
		Username: username,
		Password: password,
	}
	err = repository.CreateUser(ctx, user)
	if err != nil {
		return nil, apperror.ErrInternal.WithDetail("创建用户失败")
	}

	return &dto.UserRegisterRes{
		Username: username,
	}, nil
}

func (s *UserSVC) Login(ctx context.Context, req *dto.UserLoginReq) (*dto.UserLoginRes, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	user, err := repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, apperror.ErrInternal.WithDetail("查询用户失败")
	}
	if user == nil {
		return nil, apperror.ErrInvalidCredential
	}

	if err = mycrypt.CheckHashAndPassword(user.Password, password); err != nil {
		return nil, apperror.ErrInvalidCredential
	}

	token, err := myjwt.GenerateToken(myjwt.MyData{Username: user.Username})
	if err != nil {
		return nil, apperror.ErrInternal.WithDetail("生成 token 失败")
	}

	return &dto.UserLoginRes{
		Username: user.Username,
		Token:    token,
	}, nil
}

func (s *UserSVC) GetProfile(ctx context.Context, req *dto.UserProfileReq) (*dto.UserProfileRes, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, apperror.ErrInvalidParam.WithDetail("username 不能为空")
	}

	user, err := repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, apperror.ErrInternal.WithDetail("查询用户失败")
	}
	if user == nil {
		return nil, apperror.ErrUserNotFound
	}

	return &dto.UserProfileRes{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
