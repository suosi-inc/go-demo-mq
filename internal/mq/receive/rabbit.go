package receive

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
	queueName := config.Cfg.RabbitQueue.Simple.Name
	simple, _ := rabbits.NewSimple(queueName)

	// 启动多个协程，每个协程都是一个 Consumer
	goroutines := 1
	if config.Cfg.RabbitQueue.Simple.Goroutines > 0 {
		goroutines = config.Cfg.RabbitQueue.Simple.Goroutines
	}
	for i := 0; i < goroutines; i++ {
		no := i
		go func() {
			if msgs, err := simple.Receive(); err == nil {
				for msg := range msgs {
					var demo domain.Demo
					if e := json.Unmarshal(msg.Body, &demo); e == nil {
						log.Infof("Received-%d simple message : %+v", no, demo)

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
	exchangeName := config.Cfg.RabbitQueue.Topic.Exchange
	queueName := config.Cfg.RabbitQueue.Topic.Name
	routingKeys := fun.SliceTrim(config.Cfg.RabbitQueue.Topic.RoutingKeys)

	topic, _ := rabbits.NewTopic(exchangeName)

	// 启动多个协程, 每个协程都是一个 Consumer
	goroutines := 1
	if config.Cfg.RabbitQueue.Topic.Goroutines > 0 {
		goroutines = config.Cfg.RabbitQueue.Topic.Goroutines
	}
	for i := 0; i < goroutines; i++ {
		no := i
		go func() {
			if msgs, err := topic.ReceiveWithRoutingKeys(queueName, routingKeys); err == nil {
				for msg := range msgs {
					var demo domain.Demo
					if e := json.Unmarshal(msg.Body, &demo); e == nil {
						log.Infof("Received-%d topic message : %+v", no, demo)

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
