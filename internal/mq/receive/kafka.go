package receive

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
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
	groupId := config.Cfg.KafkaQueue.GroupId

	r := kafkas.NewReader(topic, groupId)
	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {
			log.Fatalf("Failed to close Reader", err.Error())
		}
	}(r)

	ctx, cancel := context.WithCancel(context.Background())

	// 启动多个协程消费数据, Kafka Reader 根据 MaxBytes 获取批量数据
	goroutines := 1
	if config.Cfg.KafkaQueue.Goroutines > 0 {
		goroutines = config.Cfg.KafkaQueue.Goroutines
	}
	for i := 0; i < goroutines; i++ {
		no := i
		go func() {
			for {
				// FetchMessage 需要手动 commit, 自动使用 ReadMessage
				if m, err := r.FetchMessage(ctx); err == nil {
					var demo domain.Demo
					if e := json.Unmarshal(m.Value, &demo); e == nil {
						log.Info("Received topic message:", log.Int("no", no), log.Any("demo", demo))

						// 模拟耗时
						s := fun.RandomInt(300, 800)
						time.Sleep(time.Millisecond * time.Duration(s))

					}

					// commit
					if e := r.CommitMessages(ctx, m); e != nil {
						log.Errorf("Failed commit message", no, err.Error())
					}
				} else {
					log.Error("Received topic error")
					break
				}
			}

			if err := r.Close(); err != nil {
				log.Fatal("Failed to close Reader")
			}
		}()
	}

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown and close kafka
	cancel()
}
