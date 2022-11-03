package mq

import (
	"github.com/suosi-inc/go-demo/mq/internal/mq/receive"
	"github.com/suosi-inc/go-demo/mq/internal/mq/send"
	"github.com/suosi-inc/go-demo/mq/internal/mq/setup"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

func newRun(args []string) {
	arg := args[0]
	switch arg {
	case "rabbit-send-simple":
		setup.Rabbit()
		send.RabbitSimple()
	case "rabbit-receive-simple":
		setup.Rabbit()
		receive.RabbitSimple()
	case "rabbit-send-topic":
		setup.Rabbit()
		send.RabbitTopic()
	case "rabbit-receive-topic":
		setup.Rabbit()
		receive.RabbitTopic()
	case "kafka-send-topic":
		send.KafkaTopic()
	case "kafka-receive-topic":
		receive.KafkaTopic()
	default:
		log.Error("Run args failed")
	}
}
