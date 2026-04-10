package rabbitmq_test

import (
	"context"
	"encoding/json"
	"testing"

	"finalai/internal/common/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestToJsonBytes_Success(t *testing.T) {
	msg := &rabbitmq.Message{
		SessionID: "s-1",
		Username:  "alice",
		Content:   "hello",
		IsUser:    true,
	}

	b, err := rabbitmq.ToJsonBytes(msg)
	if err != nil {
		t.Fatalf("ToJsonBytes() error = %v", err)
	}

	var got rabbitmq.Message
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got != *msg {
		t.Fatalf("unexpected message after marshal/unmarshal: got=%+v want=%+v", got, *msg)
	}
}

func TestHandleMessage_InvalidJSON_ReturnsError(t *testing.T) {
	delivery := &amqp.Delivery{Body: []byte("{invalid json")}

	err := rabbitmq.HandleMessage(context.Background(), delivery)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
