package repository

import (
	"context"
	"finalai/internal/common/mysql"
	"finalai/internal/model"
)

func GetMessagesBySessionID(sessionID string) ([]*model.Message, error) {
	var msgs []*model.Message
	err := mysql.DB.Where("session_id = ?", sessionID).Order("created_at asc").Find(&msgs).Error
	return msgs, err
}

func GetMessagesBySessionIDs(sessionIDs []string) ([]*model.Message, error) {
	var msgs []*model.Message
	if len(sessionIDs) == 0 {
		return msgs, nil
	}
	err := mysql.DB.Where("session_id IN ?", sessionIDs).Order("created_at asc").Find(&msgs).Error
	return msgs, err
}

func CreateMessage(ctx context.Context, msg *model.Message) error {
	return mysql.DB.WithContext(ctx).Create(msg).Error
}

func GetAllMessages() ([]*model.Message, error) {
	var messages []*model.Message
	err := mysql.DB.Order("created_at asc").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
