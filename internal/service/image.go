package service

import (
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/common/image"
	"io"
	"log/slog"
)

type ImageSVC struct{}

func NewImageSVC() *ImageSVC {
	return &ImageSVC{}
}

func (s *ImageSVC) RecognizeImage(req *dto.ImageRecognizeReq) (*dto.ImageRecognizeRes, error) {
	if req == nil || req.File == nil {
		return nil, apperror.ErrInvalidParam.WithDetail("image 不能为空")
	}

	modelPath := "/root/models/mobilenetv2/mobilenetv2-7.onnx"
	labelPath := "/root/imagenet_classes.txt"
	inputH, inputW := 224, 224

	recognizer, err := image.NewImageRecognizer(modelPath, labelPath, inputH, inputW)
	if err != nil {
		slog.Error("Failed to create ImageRecognizer: " + err.Error())
		return nil, apperror.ErrServerBusy.WithCause(err)
	}
	defer recognizer.Close()

	src, err := req.File.Open()
	if err != nil {
		slog.Error("Failed to open image file: " + err.Error())
		return nil, apperror.ErrInvalidParam.WithCause(err)
	}
	defer src.Close()

	buf, err := io.ReadAll(src)
	if err != nil {
		slog.Error("Failed to read image file: " + err.Error())
		return nil, apperror.ErrInternal.WithCause(err)
	}

	className, err := recognizer.PredictFromBuffer(buf)
	if err != nil {
		return nil, apperror.ErrServerBusy.WithCause(err)
	}

	return &dto.ImageRecognizeRes{ClassName: className}, nil
}
