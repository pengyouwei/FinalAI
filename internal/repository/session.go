package repository

import (
	"errors"
	"finalai/internal/common/mysql"
	"finalai/internal/model"

	"gorm.io/gorm"
)

func GetSessionsByUserName(username string) ([]*model.Session, error) {
	var sessions []*model.Session
	err := mysql.DB.Where("username = ?", username).Order("created_at desc").Find(&sessions).Error
	return sessions, err
}

func CreateSession(session *model.Session) (*model.Session, error) {
	err := mysql.DB.Create(session).Error
	return session, err
}

func GetSessionByID(sessionID string) (*model.Session, error) {
	var session model.Session
	err := mysql.DB.Where("id = ?", sessionID).First(&session).Error
	return &session, err
}

func DeleteSessionHistory(username, sessionID string) (bool, error) {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("session_id = ? AND username = ?", sessionID, username).Delete(&model.Message{}).Error; err != nil {
			return err
		}

		res := tx.Where("id = ? AND username = ?", sessionID, username).Delete(&model.Session{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
