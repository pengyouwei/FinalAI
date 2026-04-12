package dto

import (
	"mime/multipart"
	"net/http"
)

type UserRegisterReq struct {
	Username        string `json:"username" validate:"required,min=3,max=20"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserLoginReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type UserProfileReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
}

type SessionListReq struct {
	Username string `json:"-" validate:"required,min=3,max=20"`
}

type CreateSessionAndSendReq struct {
	Username     string `json:"-" validate:"required,min=3,max=20"`
	UserQuestion string `json:"question" validate:"required"`
	ModelType    string `json:"modelType" validate:"required"`
}

type CreateStreamSessionReq struct {
	Username     string `json:"-" validate:"required,min=3,max=20"`
	UserQuestion string `json:"question" validate:"required"`
	ModelType    string `json:"modelType" validate:"required"`
}

type StreamMessageReq struct {
	Username     string              `json:"-" validate:"required,min=3,max=20"`
	SessionID    string              `json:"sessionId" validate:"required"`
	UserQuestion string              `json:"question" validate:"required"`
	ModelType    string              `json:"modelType" validate:"required"`
	Writer       http.ResponseWriter `json:"-"`
}

type CreateStreamSessionAndSendReq struct {
	Username     string              `json:"-" validate:"required,min=3,max=20"`
	UserQuestion string              `json:"question" validate:"required"`
	ModelType    string              `json:"modelType" validate:"required"`
	Writer       http.ResponseWriter `json:"-"`
}

type ChatSendReq struct {
	Username     string `json:"-" validate:"required,min=3,max=20"`
	SessionID    string `json:"sessionId" validate:"required"`
	UserQuestion string `json:"question" validate:"required"`
	ModelType    string `json:"modelType" validate:"required"`
}

type ChatHistoryReq struct {
	Username  string `json:"-" validate:"required,min=3,max=20"`
	SessionID string `json:"sessionId" validate:"required"`
}

type DeleteSessionReq struct {
	Username  string `json:"-" validate:"required,min=3,max=20"`
	SessionID string `json:"sessionId" validate:"required"`
}

type ChatStreamSendReq struct {
	Username     string              `json:"-" validate:"required,min=3,max=20"`
	SessionID    string              `json:"sessionId" validate:"required"`
	UserQuestion string              `json:"question" validate:"required"`
	ModelType    string              `json:"modelType" validate:"required"`
	Writer       http.ResponseWriter `json:"-"`
}

type ImageRecognizeReq struct {
	File *multipart.FileHeader `json:"-"`
}

type UploadRagFileReq struct {
	Username string                `json:"-" validate:"required,min=3,max=20"`
	File     *multipart.FileHeader `json:"-" validate:"required"`
}
