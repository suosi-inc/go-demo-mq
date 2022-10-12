package rabbit

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/di"
)

type topic struct {
	Exchange   string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
}

func NewTopic(exchange string) (*topic, error) {
	return NewTopicWithOptions(exchange, true, false, false)
}

func NewTopicWithOptions(exchange string, durable bool, autoDelete bool, exclusive bool) (*topic, error) {
	channel := di.GetRabbit()

	c := &topic{
		Exchange:   exchange,
		Durable:    durable,
		AutoDelete: autoDelete,
		Exclusive:  exclusive,
	}

	err := channel.ExchangeDeclare(
		exchange,
		"topic",
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

func (q *topic) Send(routingKey string, body []byte) error {
	channel := di.GetRabbit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx,
		q.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	return err
}

func (q *topic) Receive(queueName string, routingKeys []string) (<-chan amqp.Delivery, error) {
	channel := di.GetRabbit()

	queue, err := channel.QueueDeclare(
		queueName,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		false,
		nil,
	)

	for _, routingKey := range routingKeys {
		channel.QueueBind(
			queue.Name,
			routingKey,
			q.Exchange,
			false,
			nil,
		)
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	return msgs, err
}
