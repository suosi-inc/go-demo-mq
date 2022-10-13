package receiver

import (
	"context"
	"encoding/json"
	"time"

	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/suosi-inc/go-demo/mq/internal/mq/data/domain"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/kafkas"
	"github.com/suosi-inc/go-demo/mq/internal/pkg/log"
	"github.com/x-funs/go-fun"
)

func KafkaTopic() {
	// goroutines
	// goroutines := 1
	// if config.Cfg.KafkaQueue.Goroutines > 0 {
	// 	goroutines = config.Cfg.KafkaQueue.Goroutines
	// }

	topic := config.Cfg.KafkaQueue.Topic
	groupId := config.Cfg.KafkaQueue.GroupId

	// for i := 0; i < goroutines; i++ {
	// 	no := i
	// 	go func() {
	ctx := context.Background()
	r := kafkas.NewReader(topic, groupId)
	// r.SetOffset(0)

	for {
		if m, err := r.FetchMessage(ctx); err == nil {
			var demo domain.Demo
			if e := json.Unmarshal(m.Value, &demo); e == nil {
				log.Info("Received topic message:", log.Int("no", 1), log.Any("demo", demo))

				// 模拟耗时
				s := fun.RandomInt(300, 800)
				time.Sleep(time.Millisecond * time.Duration(s))

				_ = r.CommitMessages(ctx, m)
			}
		} else {
			log.Error("Received topic error")
			break
		}
	}
	// }()
	// }
}
