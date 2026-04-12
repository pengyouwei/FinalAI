package controller

import (
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/controller/response"
	jwtauth "finalai/internal/middleware/jwt"
	"finalai/internal/service"
	"strings"

	"github.com/labstack/echo/v5"
)

type FileHandler struct {
	fileSVC *service.FileSVC
}

func NewFileHandler() *FileHandler {
	return &FileHandler{fileSVC: service.NewFileSVC()}
}

func (h *FileHandler) UploadRagFile(c *echo.Context) error {
	username, ok := c.Get(jwtauth.ContextKeyUsername).(string)
	if !ok || strings.TrimSpace(username) == "" {
		return response.Error(c, apperror.ErrUnauthorized.WithMessage("未登录或 token 无效"))
	}

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, apperror.ErrInvalidParam.WithDetail(err.Error()))
	}

	res, appErr := h.fileSVC.UploadRagFile(&dto.UploadRagFileReq{
		Username: username,
		File:     uploadedFile,
	})
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return response.Success(c, "RAG 文件上传成功", res)
}
