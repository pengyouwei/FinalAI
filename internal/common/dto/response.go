package dto

import (
	"finalai/internal/model"
	"time"
)

type UserRegisterRes struct {
	Username string `json:"username"`
}

type UserLoginRes struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserProfileRes struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SessionListRes struct {
	Sessions []model.SessionInfo `json:"sessions"`
}

type CreateSessionAndSendRes struct {
	SessionID     string `json:"sessionId"`
	AiInformation string `json:"information"`
}

type CreateStreamSessionRes struct {
	SessionID string `json:"sessionId"`
}

type StreamDoneRes struct {
	Done bool `json:"done"`
}

type ChatSendRes struct {
	AiInformation string `json:"information"`
}

type ChatHistoryRes struct {
	History []model.History `json:"history"`
}

type ImageRecognizeRes struct {
	ClassName string `json:"class_name"`
}
