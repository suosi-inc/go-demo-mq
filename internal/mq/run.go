package mq

import (
	"github.com/suosi-inc/go-demo/mq/internal/mq/receive"
	"github.com/suosi-inc/go-demo/mq/internal/mq/send"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

func newRun(args []string) {
	arg := args[0]
	switch arg {
	case "rabbit-send-simple":
		send.RabbitSimple()
	case "rabbit-receive-simple":
		receive.RabbitSimple()
	case "rabbit-send-topic":
		send.RabbitTopic()
	case "rabbit-receive-topic":
		receive.RabbitTopic()
	case "kafka-send-topic":
		send.KafkaTopic()
	case "kafka-receive-topic":
		receive.KafkaTopic()
	default:
		log.Error("Run args failed")
	}
}
