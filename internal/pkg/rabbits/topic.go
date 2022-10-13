package rabbits

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

// NewTopic 声明 Topic 队列
func NewTopic(exchange string) (*topic, error) {
	return NewTopicWithOptions(exchange, true, false, false)
}

// NewTopicWithOptions 声明 Topic 队列
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

// Send 发送消息
func (q *topic) Send(body []byte, routingKey string) error {
	return q.SendWithOptions(body, routingKey, 2, "text/plain")
}

// SendWithMode 发送消息
func (q *topic) SendWithMode(body []byte, routingKey string, mode uint8) error {
	return q.SendWithOptions(body, routingKey, mode, "text/plain")
}

// SendWithOptions 发送消息
func (q *topic) SendWithOptions(body []byte, routingKey string, mode uint8, contentType string) error {
	channel := di.GetRabbit()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx,
		q.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Body:         body,
			DeliveryMode: mode,
			ContentType:  contentType,
		},
	)

	return err
}

// Receive 接收消息
func (q *topic) Receive(queueName string) (<-chan amqp.Delivery, error) {
	routingKeys := []string{"#"}
	return q.ReceiveWithRoutingKeys(queueName, routingKeys)
}

// ReceiveWithRoutingKeys 接收消息
func (q *topic) ReceiveWithRoutingKeys(queueName string, routingKeys []string) (<-chan amqp.Delivery, error) {
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
