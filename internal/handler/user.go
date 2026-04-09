package handler

import (
	"context"
	"finalai/internal/dto"
	"finalai/internal/service"
	"net/http"
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
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "Invalid request body: " + err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "Validation failed: " + err.Error(),
		})
	}

	res, err := h.userSVC.Register(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  "Failed to register user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "用户注册成功",
		"data": res, // 指针类型也可以直接返回
	})
}

func (h *UserHandler) Login(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := new(dto.UserLoginReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "Invalid request body: " + err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"code": -1,
			"msg":  "Validation failed: " + err.Error(),
		})
	}

	res, err := h.userSVC.Login(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid username or password") {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]any{
			"code": -1,
			"msg":  "Failed to login: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "用户登录成功",
		"data": res,
	})
}
