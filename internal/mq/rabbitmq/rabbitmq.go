package rabbitmq

import (
	"finalai/internal/config"
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

func Init() {
	config := config.GetConfig().RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Vhost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		panic("Failed to connect to [RabbitMQ]: " + err.Error())
	}
	Conn = conn

	slog.Info("Successfully connected to [RabbitMQ]")
}

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	Exchange string
	Key      string
}

func NewRabbitMQ(exchange, key string) *RabbitMQ {
	return &RabbitMQ{
		Exchange: exchange,
		Key:      key,
	}
}

func (r *RabbitMQ) Destroy() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func NewWorkerRabbitMQ(queue string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", queue)
	rabbitmq.conn = Conn

	var err error
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		panic("Failed to create channel: " + err.Error())
	}

	return rabbitmq
}

func (r *RabbitMQ) Publish(body []byte) error {
	q, err := r.channel.QueueDeclare(r.Key, false, false, false, false, nil)
	if err != nil {
		return err
	}
	return r.channel.Publish(
		r.Exchange,
		q.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consume(handle func(msg *amqp.Delivery) error) {
	// 创建队列
	q, err := r.channel.QueueDeclare(r.Key, false, false, false, false, nil)
	if err != nil {
		panic("Failed to declare queue: " + err.Error())
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		panic("Failed to register consumer: " + err.Error())
	}

	for msg := range msgs {
		if err := handle(&msg); err != nil {
			slog.Error("Failed to handle message: " + err.Error())
		}
	}
}
