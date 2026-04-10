package response

import (
	"finalai/internal/common/apperror"
	"net/http"

	"github.com/labstack/echo/v5"
)

type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func Success(c *echo.Context, msg string, data any) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code: apperror.CodeOK,
		Msg:  msg,
		Data: data,
	})
}

func Error(c *echo.Context, err *apperror.Error) error {
	if err == nil {
		err = apperror.ErrInternal
	}

	return c.JSON(err.HTTPStatus, APIResponse{
		Code: err.Code,
		Msg:  err.Message,
	})
}

func ErrorFrom(c *echo.Context, err error) error {
	if appErr := apperror.As(err); appErr != nil {
		return Error(c, appErr)
	}
	return Error(c, apperror.ErrInternal)
}
