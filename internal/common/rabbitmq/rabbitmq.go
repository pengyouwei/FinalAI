package rabbitmq

import (
	"finalai/internal/config"
	"fmt"
	"log/slog"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
)

const (
	MessageQueueName = "message.queue"
)

func Init() {
	config := config.GetConfig().RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Vhost,
	)

	c, err := amqp.Dial(url)
	if err != nil {
		panic("Failed to connect to [RabbitMQ]: " + err.Error())
	}
	conn = c

	slog.Info("Successfully connected to [RabbitMQ]")
}

func CloseConn() {
	if conn != nil {
		_ = conn.Close()
	}
	slog.Info("Successfully closed [RabbitMQ] connection")
}

func declareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // 队列名称
		false,     // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

// Publish 每次发布都临时创建 channel，发布后立即关闭。
func Publish(queueName string, body []byte) error {
	if conn == nil {
		return fmt.Errorf("rabbitmq connection is not initialized")
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareQueue(ch, queueName)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

type Consumer struct {
	channel   *amqp.Channel
	closeOnce sync.Once
}

// StartConsumer 启动一个长连接消费者，生命周期由 Consumer.Close 控制。
func StartConsumer(queueName string, handle func(msg *amqp.Delivery) error) *Consumer {
	if conn == nil {
		panic("rabbitmq connection is not initialized")
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to create channel: " + err.Error())
	}

	q, err := declareQueue(ch, queueName)
	if err != nil {
		_ = ch.Close()
		panic("Failed to declare queue: " + err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		_ = ch.Close()
		panic("Failed to register consumer: " + err.Error())
	}

	c := &Consumer{channel: ch}
	go func() {
		for msg := range msgs {
			if err := handle(&msg); err != nil {
				slog.Error("Failed to handle message: " + err.Error())
			}
		}
	}()

	return c
}

func (c *Consumer) Close() {
	if c == nil {
		return
	}

	c.closeOnce.Do(func() {
		if c.channel != nil {
			_ = c.channel.Close()
		}
	})
}
