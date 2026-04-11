package rabbitmq

import (
	"context"
	"encoding/json"
	"finalai/internal/model"
	"finalai/internal/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	IsUser    bool   `json:"is_user"`
}

func ToJsonBytes(msg *Message) ([]byte, error) {
	return json.Marshal(msg)
}

func HandleMessage(ctx context.Context, delivery *amqp.Delivery) error {
	msg := &Message{}
	if err := json.Unmarshal(delivery.Body, msg); err != nil {
		return err
	}
	return repository.CreateMessage(ctx, &model.Message{
		SessionID: msg.SessionID,
		Username:  msg.Username,
		Content:   msg.Content,
		IsUser:    msg.IsUser,
	})
}
