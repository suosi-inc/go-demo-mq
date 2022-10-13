package send

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/mq/data/domain"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/kafkas"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
	"github.com/x-funs/go-fun"
)

func KafkaTopic() {
	topic := config.Cfg.KafkaQueue.Topic
	w := kafkas.NewWriter(topic)

	var id int

	ctx := context.Background()

	for {
		id++
		msg := domain.Demo{
			Id:   id,
			Name: fun.RandomLetter(4),
			Time: fun.Date(fun.DatetimeMilliPattern),
		}

		msgJson := fun.ToJson(msg)

		err := w.WriteMessages(ctx,
			kafka.Message{
				Value: fun.Bytes(msgJson),
			},
		)

		if err == nil {
			log.Info("Send topic success", log.String("msg", msgJson))
		} else {
			log.Error("Send topic error")
		}

		// sleep
		time.Sleep(time.Millisecond * 1000)
	}
}
