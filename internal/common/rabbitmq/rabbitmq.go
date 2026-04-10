package rabbitmq

import (
	"finalai/internal/config"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

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

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to create channel: " + err.Error())
	}
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
}

func CloseConn() {
	if conn != nil {
		_ = conn.Close()
	}
}

func (r *RabbitMQ) Publish(queueName string, body []byte) error {
	q, err := r.declareQueue(queueName)
	if err != nil {
		return err
	}
	return r.channel.Publish(
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

func (r *RabbitMQ) Consume(queueName string, handle func(msg *amqp.Delivery) error) {
	// 创建队列
	q, err := r.declareQueue(queueName)
	if err != nil {
		panic("Failed to declare queue: " + err.Error())
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		panic("Failed to register consumer: " + err.Error())
	}

	go func() {
		for msg := range msgs {
			if err := handle(&msg); err != nil {
				slog.Error("Failed to handle message: " + err.Error())
			}
		}
	}()
}

func (r *RabbitMQ) declareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName, // 队列名称
		false,     // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}
