package mq

import (
	"github.com/suosi-inc/go-demo/mq/internal/mq/receiver"
	"github.com/suosi-inc/go-demo/mq/internal/mq/sender"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
)

func newRun(args []string) {
	arg := args[0]
	switch arg {
	case "rabbit-sender-simple":
		sender.RabbitSimple()
	case "rabbit-receiver-simple":
		receiver.RabbitSimple()
	case "rabbit-sender-topic":
		sender.RabbitTopic()
	case "rabbit-receiver-topic":
		receiver.RabbitTopic()
	default:
		log.Error("Run args failed")
	}
}
