package rabbits

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/di"
)

type simple struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
}

// NewSimple 声明简单队列
func NewSimple(name string) (*simple, error) {
	return NewSimpleWithOptions(name, true, false, false)
}

// NewSimpleWithOptions 声明简单队列
func NewSimpleWithOptions(name string, durable bool, autoDelete bool, exclusive bool) (*simple, error) {
	channel := di.GetRabbit()

	c := &simple{
		Name:       name,
		Durable:    durable,
		AutoDelete: autoDelete,
		Exclusive:  exclusive,
	}

	// declare
	_, err := channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

// Send 发送消息
func (q *simple) Send(body []byte) error {
	return q.SendWithOptions(body, 2, "text/plain")
}

// SendWithMode 发送消息
func (q *simple) SendWithMode(body []byte, mode uint8) error {
	return q.SendWithOptions(body, mode, "text/plain")
}

// SendWithOptions 发送消息s
func (q *simple) SendWithOptions(body []byte, mode uint8, contentType string) error {
	channel := di.GetRabbit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  contentType,
			Body:         body,
			DeliveryMode: mode,
		},
	)

	return err
}

// Receive 接收消息
func (q *simple) Receive() (<-chan amqp.Delivery, error) {
	channel := di.GetRabbit()

	msgs, err := channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	return msgs, err
}
