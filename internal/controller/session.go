package controller

import (
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/controller/response"
	jwtauth "finalai/internal/middleware/jwt"
	"finalai/internal/service"
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo/v5"
)

type SessionHandler struct {
	sessionSVC *service.SessionSVC
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{sessionSVC: service.NewSessionSVC()}
}

func (h *SessionHandler) GetUserSessionsByUserName(c *echo.Context) error {
	username, err := h.usernameFromContext(c)
	if err != nil {
		return response.Error(c, err)
	}

	res, appErr := h.sessionSVC.GetUserSessionsByUserName(&dto.SessionListReq{Username: username})
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return response.Success(c, "获取会话列表成功", res)
}

func (h *SessionHandler) CreateSessionAndSendMessage(c *echo.Context) error {
	username, err := h.usernameFromContext(c)
	if err != nil {
		return response.Error(c, err)
	}

	req := new(dto.CreateSessionAndSendReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}
	req.Username = username

	res, appErr := h.sessionSVC.CreateSessionAndSendMessage(req)
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return response.Success(c, "创建会话并发送成功", res)
}

func (h *SessionHandler) CreateStreamSessionAndSendMessage(c *echo.Context) error {
	username, err := h.usernameFromContext(c)
	if err != nil {
		return response.Error(c, err)
	}

	req := new(dto.CreateStreamSessionReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}

	raw := c.Response()
	raw.Header().Set("Content-Type", "text/event-stream")
	raw.Header().Set("Cache-Control", "no-cache")
	raw.Header().Set("Connection", "keep-alive")
	raw.Header().Set("Access-Control-Allow-Origin", "*")
	raw.Header().Set("X-Accel-Buffering", "no")

	sessionRes, appErr := h.sessionSVC.CreateStreamSessionOnly(&dto.CreateStreamSessionReq{
		Username:     username,
		UserQuestion: req.UserQuestion,
	})
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	if _, err := raw.Write([]byte(fmt.Sprintf("data: {\"sessionId\": \"%s\"}\n\n", sessionRes.SessionID))); err != nil {
		return response.Error(c, apperror.ErrInternal.WithCause(err))
	}
	if flusher, ok := raw.(interface{ Flush() }); ok {
		flusher.Flush()
	}

	_, appErr = h.sessionSVC.StreamMessageToExistingSession(&dto.StreamMessageReq{
		Username:     username,
		SessionID:    sessionRes.SessionID,
		UserQuestion: req.UserQuestion,
		ModelType:    req.ModelType,
		Writer:       raw,
	})
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return nil
}

func (h *SessionHandler) ChatSend(c *echo.Context) error {
	username, err := h.usernameFromContext(c)
	if err != nil {
		return response.Error(c, err)
	}

	req := new(dto.ChatSendReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}
	req.Username = username

	res, appErr := h.sessionSVC.ChatSend(req)
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return response.Success(c, "发送消息成功", res)
}

func (h *SessionHandler) ChatStreamSend(c *echo.Context) error {
	username, err := h.usernameFromContext(c)
	if err != nil {
		return response.Error(c, err)
	}

	req := new(dto.ChatStreamSendReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}
	req.Username = username
	req.Writer = c.Response()

	raw := c.Response()
	raw.Header().Set("Content-Type", "text/event-stream")
	raw.Header().Set("Cache-Control", "no-cache")
	raw.Header().Set("Connection", "keep-alive")
	raw.Header().Set("Access-Control-Allow-Origin", "*")
	raw.Header().Set("X-Accel-Buffering", "no")

	_, appErr := h.sessionSVC.ChatStreamSend(req)
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return nil
}

func (h *SessionHandler) ChatHistory(c *echo.Context) error {
	username, err := h.usernameFromContext(c)
	if err != nil {
		return response.Error(c, err)
	}

	req := new(dto.ChatHistoryReq)
	if err := h.bindAndValidate(c, req); err != nil {
		return err
	}
	req.Username = username

	res, appErr := h.sessionSVC.GetChatHistory(req)
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return response.Success(c, "获取聊天历史成功", res)
}

func (h *SessionHandler) bindAndValidate(c *echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return response.Error(c, apperror.ErrInvalidParam.WithDetail(err.Error()))
	}
	h.injectUsernameFromContext(c, req)
	if err := c.Validate(req); err != nil {
		return response.Error(c, apperror.ErrInvalidParam.WithDetail(err.Error()))
	}
	return nil
}

func (h *SessionHandler) injectUsernameFromContext(c *echo.Context, req any) {
	username, ok := c.Get(jwtauth.ContextKeyUsername).(string)
	if !ok || strings.TrimSpace(username) == "" {
		return
	}

	rv := reflect.ValueOf(req)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return
	}

	ev := rv.Elem()
	if ev.Kind() != reflect.Struct {
		return
	}

	field := ev.FieldByName("Username")
	if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.String {
		return
	}

	if strings.TrimSpace(field.String()) == "" {
		field.SetString(username)
	}
}

func (h *SessionHandler) usernameFromContext(c *echo.Context) (string, *apperror.Error) {
	username, ok := c.Get(jwtauth.ContextKeyUsername).(string)
	if !ok || strings.TrimSpace(username) == "" {
		return "", apperror.ErrUnauthorized.WithMessage("未登录或 token 无效")
	}
	return username, nil
}
