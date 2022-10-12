package sender

import (
	"time"

	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/mq/data/domain"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
	rabbit2 "github.com/suosi-inc/go-demo/mq/internal/pkg/rabbit"
	"github.com/x-funs/go-fun"
)

func RabbitSimple() {
	queueName := config.Cfg.Queue.Simple.Name
	if simple, err := rabbit2.NewSimple(queueName); err == nil {
		var id int

		for {
			id++
			msg := domain.Demo{
				Id:   id,
				Name: fun.RandomLetter(4),
				Time: fun.Date(fun.DatetimeMilliPattern),
			}

			msgJson := fun.ToJson(msg)

			if err := simple.Send(fun.Bytes(msgJson)); err == nil {
				log.Info("Send simple success", log.String("msg", msgJson))
			} else {
				log.Error("Send simple error")
			}

			// sleep
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func RabbitTopic() {
	exchangeName := config.Cfg.Queue.Topic.Exchange
	if topic, err := rabbit2.NewTopic(exchangeName); err == nil {

		var id int
		var routingKey string

		for {
			id++

			// routingKey
			if id%2 == 0 {
				routingKey = "log.info"
			} else {
				routingKey = "log.warn"
			}

			msg := domain.Demo{
				Id:   id,
				Name: fun.RandomLetter(4),
				Time: fun.Date(fun.DatetimeMilliPattern),
				Key:  routingKey,
			}

			msgJson := fun.ToJson(msg)

			if err := topic.Send(routingKey, fun.Bytes(msgJson)); err == nil {
				log.Info("Send topic success", log.String("key", routingKey), log.String("msg", msgJson))
			} else {
				log.Error("Send topic error")
			}

			// sleep
			time.Sleep(time.Millisecond * 10)
		}
	}
}
