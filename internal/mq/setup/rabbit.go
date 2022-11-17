package setup

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/di"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

func Rabbit() {
	if connection, err := amqp.Dial(config.Cfg.Rabbit.Amqp); err == nil {
		if channel, err := connection.Channel(); err == nil {

			// Qos prefetch 仅仅在手动 ack 时有效
			if config.Cfg.Rabbit.Prefetch > 0 {
				_ = channel.Qos(config.Cfg.Rabbit.Prefetch, 0, false)
			}

			di.SetRabbit(channel)

			log.Info("Setup rabbit success")
		} else {
			log.Fatal("Setup rabbit channel failed", log.String("err", err.Error()))
		}
	} else {
		log.Fatal("Setup rabbit connection failed", log.String("err", err.Error()))
	}
}
