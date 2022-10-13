package mq

import (
	"github.com/suosi-inc/go-demo/mq/internal/mq/receiver"
	"github.com/suosi-inc/go-demo/mq/internal/mq/sender"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

func newRun(args []string) {
	arg := args[0]
	switch arg {
	case "rabbit-send-simple":
		sender.RabbitSimple()
	case "rabbit-receive-simple":
		receiver.RabbitSimple()
	case "rabbit-send-topic":
		sender.RabbitTopic()
	case "rabbit-receive-topic":
		receiver.RabbitTopic()
	case "kafka-send-topic":
		sender.KafkaTopic()
	case "kafka-receive-topic":
		receiver.KafkaTopic()
	default:
		log.Error("Run args failed")
	}
}
