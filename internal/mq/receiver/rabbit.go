package receiver

import (
	"encoding/json"
	"time"

	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/mq/data/domain"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/rabbits"
	"github.com/x-funs/go-fun"
)

func RabbitSimple() {
	// goroutines
	goroutines := 1
	if config.Cfg.RabbitQueue.Simple.Goroutines > 0 {
		goroutines = config.Cfg.RabbitQueue.Simple.Goroutines
	}

	queueName := config.Cfg.RabbitQueue.Simple.Name
	simple, _ := rabbits.NewSimple(queueName)

	for i := 0; i < goroutines; i++ {
		no := i
		go func() {
			if msgs, err := simple.Receive(); err == nil {
				for msg := range msgs {
					var demo domain.Demo
					if e := json.Unmarshal(msg.Body, &demo); e == nil {
						log.Info("Received simple message:", log.Int("no", no), log.Any("demo", demo))

						// 模拟耗时
						s := fun.RandomInt(300, 800)
						time.Sleep(time.Millisecond * time.Duration(s))

						// 消息确认
						_ = msg.Ack(false)
					}
				}
			} else {
				log.Error("Received topic error")
			}
		}()
	}
}

func RabbitTopic() {
	// goroutines
	goroutines := 1
	if config.Cfg.RabbitQueue.Topic.Goroutines > 0 {
		goroutines = config.Cfg.RabbitQueue.Topic.Goroutines
	}

	exchangeName := config.Cfg.RabbitQueue.Topic.Exchange
	queueName := config.Cfg.RabbitQueue.Topic.Name
	routingKeys := fun.SliceTrim(config.Cfg.RabbitQueue.Topic.RoutingKeys)

	topic, _ := rabbits.NewTopic(exchangeName)

	for i := 0; i < goroutines; i++ {
		no := i
		go func() {
			if msgs, err := topic.ReceiveWithRoutingKeys(queueName, routingKeys); err == nil {
				for msg := range msgs {
					var demo domain.Demo
					if e := json.Unmarshal(msg.Body, &demo); e == nil {
						log.Info("Received topic message:", log.Int("no", no), log.Any("demo", demo))

						// 模拟耗时
						s := fun.RandomInt(300, 800)
						time.Sleep(time.Millisecond * time.Duration(s))

						// 消息确认
						_ = msg.Ack(false)
					}
				}
			} else {
				log.Error("Received topic error")
			}
		}()
	}
}
