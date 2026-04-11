package controller

import (
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/controller/response"
	"finalai/internal/service"

	"github.com/labstack/echo/v5"
)

type ImageHandler struct {
	imageSVC *service.ImageSVC
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{imageSVC: service.NewImageSVC()}
}

func (h *ImageHandler) RecognizeImage(c *echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return response.Error(c, apperror.ErrInvalidParam.WithDetail(err.Error()))
	}

	res, appErr := h.imageSVC.RecognizeImage(&dto.ImageRecognizeReq{File: file})
	if appErr != nil {
		return response.ErrorFrom(c, appErr)
	}

	return response.Success(c, "图片识别成功", res)
}
