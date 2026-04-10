package repository

import (
	"context"
	"finalai/internal/common/mysql"
	"finalai/internal/model"
)

func SaveMessage(ctx context.Context, msg *model.Message) error {
	return mysql.DB.WithContext(ctx).Create(msg).Error
}
