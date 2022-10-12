package rabbit

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

func NewSimple(name string) (*simple, error) {
	return NewSimpleWithOptions(name, true, false, false)
}

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

func (q *simple) Send(body []byte) error {
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
			ContentType: "text/plain",
			Body:        body,
		},
	)

	return err
}

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
