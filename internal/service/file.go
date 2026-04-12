package service

import (
	"context"
	"finalai/internal/common/apperror"
	"finalai/internal/common/dto"
	"finalai/internal/common/rag"
	"finalai/internal/config"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileSVC struct{}

func NewFileSVC() *FileSVC {
	return &FileSVC{}
}

// 上传rag相关文件（这里只允许文本文件）
// 其实可以直接将其向量化进行保存，但这边依旧存储到服务器上以便后续可以在服务器上查看历史 RAG 文件。
func (s *FileSVC) UploadRagFile(req *dto.UploadRagFileReq) (*dto.UploadRagFileRes, error) {
	if req == nil || req.File == nil {
		return nil, apperror.ErrInvalidParam.WithDetail("file 不能为空")
	}

	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, apperror.ErrUnauthorized.WithMessage("未登录或 token 无效")
	}

	if err := validateRagFile(req.File); err != nil {
		return nil, err
	}

	// 创建用户目录
	userDir := filepath.Join("uploads", username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		slog.Error("Failed to create user directory", "dir", userDir, "error", err)
		return nil, apperror.ErrInternal.WithCause(err)
	}

	oldFiles, err := listFilesInDir(userDir)
	if err != nil {
		slog.Error("Failed to list files in user directory", "dir", userDir, "error", err)
		return nil, apperror.ErrInternal.WithCause(err)
	}

	if err := removeAllFilesInDir(userDir); err != nil {
		slog.Error("Failed to clean user directory", "dir", userDir, "error", err)
		return nil, apperror.ErrInternal.WithCause(err)
	}

	for _, oldFile := range oldFiles {
		if err := rag.DeleteIndex(context.Background(), oldFile); err != nil {
			slog.Warn("Failed to delete old rag index", "file", oldFile, "error", err)
		}
	}

	ext := strings.ToLower(filepath.Ext(req.File.Filename))
	filename := uuid.NewString() + ext
	filePath := filepath.Join(userDir, filename)

	// 打开上传的文件
	src, err := req.File.Open()
	if err != nil {
		slog.Error("Failed to open uploaded file", "error", err)
		return nil, apperror.ErrInvalidParam.WithCause(err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		slog.Error("Failed to create destination file", "path", filePath, "error", err)
		return nil, apperror.ErrInternal.WithCause(err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		slog.Error("Failed to copy file content", "error", err)
		return nil, apperror.ErrInternal.WithCause(err)
	}

	slog.Info("File uploaded successfully", "path", filePath)

	indexer, err := rag.NewRAGIndexer(filename, config.GetConfig().RagModelConfig.RagEmbeddingModel)
	if err != nil {
		slog.Error("Failed to create rag indexer", "error", err)
		_ = os.Remove(filePath)
		return nil, apperror.ErrServerBusy.WithCause(err)
	}

	if err := indexer.IndexFile(context.Background(), filePath); err != nil {
		slog.Error("Failed to index rag file", "path", filePath, "error", err)
		_ = os.Remove(filePath)
		_ = rag.DeleteIndex(context.Background(), filename)
		return nil, apperror.ErrServerBusy.WithCause(err)
	}

	slog.Info("RAG file indexed successfully", "file", filename)

	return &dto.UploadRagFileRes{
		FilePath: filePath,
		FileName: filename,
	}, nil
}

func validateRagFile(file *multipart.FileHeader) error {
	if file == nil {
		return apperror.ErrInvalidParam.WithDetail("file 不能为空")
	}

	name := strings.TrimSpace(file.Filename)
	if name == "" {
		return apperror.ErrInvalidParam.WithDetail("文件名不能为空")
	}

	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".txt" && ext != ".md" && ext != ".markdown" {
		return apperror.ErrInvalidParam.WithDetail("仅支持 txt 或 markdown 文件")
	}

	return nil
}

func removeAllFilesInDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil {
			return err
		}
	}

	return nil
}

func listFilesInDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}

	return files, nil
}
