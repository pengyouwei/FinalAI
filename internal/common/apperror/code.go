package apperror

import (
	"errors"
	"net/http"
)

// Error 定义业务错误：错误码、错误信息、HTTP 状态码。
type Error struct {
	Code       int
	Message    string
	HTTPStatus int
	Cause      error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *Error) clone() *Error {
	if e == nil {
		return nil
	}
	cp := *e
	return &cp
}

func (e *Error) WithCause(cause error) *Error {
	cp := e.clone()
	cp.Cause = cause
	return cp
}

func (e *Error) WithMessage(msg string) *Error {
	cp := e.clone()
	cp.Message = msg
	return cp
}

func (e *Error) WithDetail(detail string) *Error {
	if detail == "" {
		return e.clone()
	}
	cp := e.clone()
	cp.Message = cp.Message + ": " + detail
	return cp
}

func As(err error) *Error {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

const (
	CodeOK = 0

	CodeInvalidParam      = 10001
	CodeUnauthorized      = 10002
	CodeInvalidCredential = 10003
	CodeUserNotFound      = 10004
	CodeUserAlreadyExists = 10005
	CodeTokenMissing      = 10006
	CodeTokenFormat       = 10007
	CodeTokenInvalid      = 10008

	CodeInternal = 20000
)

var (
	ErrInvalidParam      = &Error{Code: CodeInvalidParam, Message: "请求参数错误", HTTPStatus: http.StatusBadRequest}
	ErrUnauthorized      = &Error{Code: CodeUnauthorized, Message: "未登录或无权限", HTTPStatus: http.StatusUnauthorized}
	ErrInvalidCredential = &Error{Code: CodeInvalidCredential, Message: "用户名或密码错误", HTTPStatus: http.StatusUnauthorized}
	ErrUserNotFound      = &Error{Code: CodeUserNotFound, Message: "用户不存在", HTTPStatus: http.StatusNotFound}
	ErrUserAlreadyExists = &Error{Code: CodeUserAlreadyExists, Message: "用户已存在", HTTPStatus: http.StatusConflict}
	ErrTokenMissing      = &Error{Code: CodeTokenMissing, Message: "缺少 Authorization 头", HTTPStatus: http.StatusUnauthorized}
	ErrTokenFormat       = &Error{Code: CodeTokenFormat, Message: "Authorization 格式错误，应为 Bearer <token>", HTTPStatus: http.StatusUnauthorized}
	ErrTokenInvalid      = &Error{Code: CodeTokenInvalid, Message: "无效或过期的 token", HTTPStatus: http.StatusUnauthorized}
	ErrInternal          = &Error{Code: CodeInternal, Message: "服务器内部错误", HTTPStatus: http.StatusInternalServerError}
)
