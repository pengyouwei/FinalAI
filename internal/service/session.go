package service

import (
	"context"
	"finalai/internal/common/aihelper"
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/model"
	"finalai/internal/repository"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

type SessionSVC struct{}

func NewSessionSVC() *SessionSVC {
	return &SessionSVC{}
}

var ctx = context.Background()

func (s *SessionSVC) GetUserSessionsByUserName(req *dto.SessionListReq) (*dto.SessionListRes, error) {
	sessions, err := repository.GetSessionsByUserName(req.Username)
	if err != nil {
		slog.Error("GetUserSessionsByUserName GetSessionsByUserName error", "error", err)
		return nil, apperror.ErrServerBusy.WithCause(err)
	}

	infos := make([]model.SessionInfo, 0, len(sessions))
	for _, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = session.ID
		}
		infos = append(infos, model.SessionInfo{
			SessionID: session.ID,
			Title:     title,
		})
	}

	return &dto.SessionListRes{Sessions: infos}, nil
}

func (s *SessionSVC) CreateSessionAndSendMessage(req *dto.CreateSessionAndSendReq) (*dto.CreateSessionAndSendRes, error) {
	//1：创建一个新的会话
	newSession := &model.Session{
		ID:       uuid.New().String(),
		Username: req.Username,
		Title:    req.UserQuestion,
	}
	createdSession, err := repository.CreateSession(newSession)
	if err != nil {
		slog.Error("CreateSessionAndSendMessage CreateSession error", "error", err)
		return nil, apperror.ErrServerBusy.WithCause(err)
	}

	//2：获取AIHelper并通过其管理消息
	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey": "your-api-key", // TODO: 从配置中获取
	}
	helper, err := manager.GetOrCreateAIHelper(req.Username, createdSession.ID, req.ModelType, config)
	if err != nil {
		slog.Error("CreateSessionAndSendMessage GetOrCreateAIHelper error", "error", err)
		return nil, apperror.ErrAIModelFail.WithCause(err)
	}

	//3：生成AI回复
	aiResponse, err_ := helper.GenerateResponse(ctx, req.Username, req.UserQuestion)
	if err_ != nil {
		slog.Error("CreateSessionAndSendMessage GenerateResponse error", "error", err_)
		return nil, apperror.ErrAIModelFail.WithCause(err_)
	}

	return &dto.CreateSessionAndSendRes{
		SessionID:     createdSession.ID,
		AiInformation: aiResponse.Content,
	}, nil
}

func (s *SessionSVC) CreateStreamSessionOnly(req *dto.CreateStreamSessionReq) (*dto.CreateStreamSessionRes, error) {
	newSession := &model.Session{
		ID:       uuid.New().String(),
		Username: req.Username,
		Title:    req.UserQuestion,
	}
	createdSession, err := repository.CreateSession(newSession)
	if err != nil {
		slog.Error("CreateStreamSessionOnly CreateSession error", "error", err)
		return nil, apperror.ErrServerBusy.WithCause(err)
	}
	return &dto.CreateStreamSessionRes{SessionID: createdSession.ID}, nil
}

func (s *SessionSVC) StreamMessageToExistingSession(req *dto.StreamMessageReq) (*dto.StreamDoneRes, error) {
	// 确保 writer 支持 Flush
	flusher, ok := req.Writer.(interface{ Flush() })
	if !ok {
		slog.Error("StreamMessageToExistingSession: streaming unsupported")
		return nil, apperror.ErrServerBusy
	}

	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey": "your-api-key", // TODO: 从配置中获取
	}
	helper, err := manager.GetOrCreateAIHelper(req.Username, req.SessionID, req.ModelType, config)
	if err != nil {
		slog.Error("StreamMessageToExistingSession GetOrCreateAIHelper error", "error", err)
		return nil, apperror.ErrAIModelFail.WithCause(err)
	}

	cb := func(msg string) {
		slog.Debug("SSE sending chunk", "len", len(msg))
		_, err := req.Writer.Write([]byte("data: " + msg + "\n\n"))
		if err != nil {
			slog.Error("SSE write error", "error", err)
			return
		}
		flusher.Flush()
		slog.Debug("SSE flushed")
	}

	_, err_ := helper.StreamResponse(ctx, req.Username, cb, req.UserQuestion)
	if err_ != nil {
		slog.Error("StreamMessageToExistingSession StreamResponse error", "error", err_)
		return nil, apperror.ErrAIModelFail.WithCause(err_)
	}

	_, err = req.Writer.Write([]byte("data: [DONE]\n\n"))
	if err != nil {
		slog.Error("StreamMessageToExistingSession write DONE error", "error", err)
		return nil, apperror.ErrAIModelFail.WithCause(err)
	}
	flusher.Flush()

	return &dto.StreamDoneRes{Done: true}, nil
}

func (s *SessionSVC) CreateStreamSessionAndSendMessage(req *dto.CreateStreamSessionAndSendReq) (*dto.CreateStreamSessionRes, error) {
	sessionRes, appErr := s.CreateStreamSessionOnly(&dto.CreateStreamSessionReq{
		Username:     req.Username,
		UserQuestion: req.UserQuestion,
	})
	if appErr != nil {
		return nil, appErr
	}

	_, appErr = s.StreamMessageToExistingSession(&dto.StreamMessageReq{
		Username:     req.Username,
		SessionID:    sessionRes.SessionID,
		UserQuestion: req.UserQuestion,
		ModelType:    req.ModelType,
		Writer:       req.Writer,
	})
	if appErr != nil {
		return sessionRes, appErr
	}

	return sessionRes, nil
}

func (s *SessionSVC) ChatSend(req *dto.ChatSendReq) (*dto.ChatSendRes, error) {
	//1：获取AIHelper
	manager := aihelper.GetGlobalManager()
	config := map[string]interface{}{
		"apiKey": "your-api-key", // TODO: 从配置中获取
	}
	helper, err := manager.GetOrCreateAIHelper(req.Username, req.SessionID, req.ModelType, config)
	if err != nil {
		slog.Error("ChatSend GetOrCreateAIHelper error", "error", err)
		return nil, apperror.ErrAIModelFail.WithCause(err)
	}

	//2：生成AI回复
	aiResponse, err_ := helper.GenerateResponse(ctx, req.Username, req.UserQuestion)
	if err_ != nil {
		slog.Error("ChatSend GenerateResponse error", "error", err_)
		return nil, apperror.ErrAIModelFail.WithCause(err_)
	}

	return &dto.ChatSendRes{AiInformation: aiResponse.Content}, nil
}

func (s *SessionSVC) GetChatHistory(req *dto.ChatHistoryReq) (*dto.ChatHistoryRes, error) {
	// 获取AIHelper中的消息历史
	manager := aihelper.GetGlobalManager()
	helper, exists := manager.GetAIHelper(req.Username, req.SessionID)
	if !exists {
		return nil, apperror.ErrSessionNotFound
	}

	messages := helper.GetMessages()
	history := make([]model.History, 0, len(messages))

	// 转换消息为历史格式（根据消息顺序或内容判断用户/AI消息）
	for i, msg := range messages {
		isUser := i%2 == 0
		history = append(history, model.History{
			IsUser:  isUser,
			Content: msg.Content,
		})
	}

	return &dto.ChatHistoryRes{History: history}, nil
}

func (s *SessionSVC) ChatStreamSend(req *dto.ChatStreamSendReq) (*dto.StreamDoneRes, error) {
	return s.StreamMessageToExistingSession(&dto.StreamMessageReq{
		Username:     req.Username,
		SessionID:    req.SessionID,
		UserQuestion: req.UserQuestion,
		ModelType:    req.ModelType,
		Writer:       req.Writer,
	})
}

func (s *SessionSVC) DeleteSessionHistory(req *dto.DeleteSessionReq) error {
	if req == nil || strings.TrimSpace(req.SessionID) == "" {
		return apperror.ErrInvalidParam.WithDetail("sessionId 不能为空")
	}

	deleted, err := repository.DeleteSessionHistory(req.Username, req.SessionID)
	if err != nil {
		slog.Error("DeleteSessionHistory DeleteSessionHistory error", "error", err)
		return apperror.ErrServerBusy.WithCause(err)
	}
	if !deleted {
		return apperror.ErrSessionNotFound
	}

	aihelper.GetGlobalManager().RemoveAIHelper(req.Username, req.SessionID)
	return nil
}
