package handler

import (
	"context"
	"finalai/internal/apperror"
	"finalai/internal/dto"
	"finalai/internal/response"
	"finalai/internal/service"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userSVC *service.UserSVC
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userSVC: service.NewUserSVC(),
	}
}

func (h *UserHandler) Register(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := new(dto.UserRegisterReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}

	res, err := h.userSVC.Register(ctx, req)
	if err != nil {
		return response.ErrorFrom(c, err)
	}

	return response.Success(c, "用户注册成功", res)
}

func (h *UserHandler) Login(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := new(dto.UserLoginReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}

	res, err := h.userSVC.Login(ctx, req)
	if err != nil {
		return response.ErrorFrom(c, err)
	}

	return response.Success(c, "用户登录成功", res)
}

func (h *UserHandler) GetProfile(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	username, ok := c.Get("username").(string)
	if !ok || strings.TrimSpace(username) == "" {
		return response.Error(c, apperror.ErrUnauthorized.WithMessage("未登录或 token 无效"))
	}

	req := &dto.UserProfileReq{Username: username}
	res, err := h.userSVC.GetProfile(ctx, req)
	if err != nil {
		return response.ErrorFrom(c, err)
	}

	return response.Success(c, "获取用户信息成功", res)
}

func (h *UserHandler) bindAndValidate(c *echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return response.Error(c, apperror.ErrInvalidParam.WithDetail(err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return response.Error(c, apperror.ErrInvalidParam.WithDetail(err.Error()))
	}

	return nil
}
