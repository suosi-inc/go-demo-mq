package kafkas

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/x-funs/go-fun"
)

const (
	KafkaTraceFile = "kafka-trace.log"
	KafkaInfo      = "[INFO]"
	KafkaError     = "[ERROR]"
)

func NewWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{}

	servers := config.Cfg.Kafka.Servers
	serverList := fun.SplitTrim(servers, ",")

	// 消息发送异步确认
	async := config.Cfg.Kafka.Write.Async

	// 生产者压缩算法
	compress := kafka.Lz4
	if !fun.Blank(config.Cfg.Kafka.Compress) {
		switch strings.ToLower(config.Cfg.Kafka.Compress) {
		case "gzip":
			compress = kafka.Gzip
		case "zstd":
			compress = kafka.Zstd
		case "snappy":
			compress = kafka.Snappy
		default:
			compress = kafka.Lz4
		}
	}

	// 日志开关
	logger := config.Cfg.Kafka.Logger

	// 是否 Sasl 认证
	if config.Cfg.Kafka.Sasl.Enable {
		var mechanism sasl.Mechanism
		if config.Cfg.Kafka.Sasl.Mechanism == "PLAIN" {
			mechanism = plain.Mechanism{
				Username: config.Cfg.Kafka.Sasl.Username,
				Password: config.Cfg.Kafka.Sasl.Password,
			}
		}

		sharedTransport := &kafka.Transport{
			SASL: mechanism,
		}

		w = &kafka.Writer{
			Addr:      kafka.TCP(serverList...),
			Topic:     topic,
			Balancer:  &kafka.Hash{},
			Transport: sharedTransport,
			Async:     async,
		}

		if logger {
			w.Logger = Logger{KafkaInfo}
			w.ErrorLogger = Logger{KafkaError}
		}
	} else {
		w = &kafka.Writer{
			Addr:        kafka.TCP(serverList...),
			Topic:       topic,
			Compression: compress,
			Balancer:    &kafka.Hash{},
			Async:       async,
		}

		if logger {
			w.Logger = Logger{KafkaInfo}
			w.ErrorLogger = Logger{KafkaError}
		}
	}

	return w
}

func NewReader(topic string, groupId string) *kafka.Reader {
	r := &kafka.Reader{}

	servers := config.Cfg.Kafka.Servers
	serverList := fun.SplitTrim(servers, ",")

	// 新的消费者组消息偏移(最初或最近)
	offset := kafka.LastOffset
	if !fun.Blank(config.Cfg.Kafka.Read.Offset) && "first" == strings.ToLower(config.Cfg.Kafka.Read.Offset) {
		offset = kafka.FirstOffset
	}

	// 消息最大字节数
	maxBytes := 102400
	if config.Cfg.Kafka.Read.MaxBytes > 0 {
		maxBytes = config.Cfg.Kafka.Read.MaxBytes
	}

	// ack 异步提交间隔，0 表示同步
	commitInterval := time.Duration(0)
	if config.Cfg.Kafka.Read.CommitInterval > 0 {
		commitInterval = time.Millisecond * time.Duration(config.Cfg.Kafka.Read.CommitInterval)
	}

	// 日志开关
	logger := config.Cfg.Kafka.Logger

	// 是否 Sasl 认证
	if config.Cfg.Kafka.Sasl.Enable {
		var mechanism sasl.Mechanism
		if config.Cfg.Kafka.Sasl.Mechanism == "PLAIN" {
			mechanism = plain.Mechanism{
				Username: config.Cfg.Kafka.Sasl.Username,
				Password: config.Cfg.Kafka.Sasl.Password,
			}
		}

		dialer := &kafka.Dialer{
			Timeout:       30 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		readerConfig := kafka.ReaderConfig{
			Dialer:           dialer,
			Brokers:          serverList,
			GroupID:          groupId,
			Topic:            topic,
			MaxBytes:         maxBytes,
			StartOffset:      offset,
			CommitInterval:   commitInterval,
			RebalanceTimeout: time.Second * 180,
			MaxAttempts:      5,
		}

		if logger {
			readerConfig.Logger = Logger{KafkaInfo}
			readerConfig.ErrorLogger = Logger{KafkaError}
		}

		r = kafka.NewReader(readerConfig)
	} else {
		readerConfig := kafka.ReaderConfig{
			Brokers:          serverList,
			GroupID:          groupId,
			Topic:            topic,
			MaxBytes:         maxBytes,
			StartOffset:      offset,
			CommitInterval:   commitInterval,
			RebalanceTimeout: time.Second * 180,
			MaxAttempts:      5,
		}

		if logger {
			readerConfig.Logger = Logger{KafkaInfo}
			readerConfig.ErrorLogger = Logger{KafkaError}
		}

		r = kafka.NewReader(readerConfig)
	}

	return r
}

type Logger struct {
	Level string
}

func (l Logger) Printf(msg string, args ...interface{}) {
	msg = fmt.Sprintf("%s [%s] %s\n", l.Level, time.Now().Format(time.RFC3339Nano), msg)
	fmt.Printf(msg, args...)

	m := fmt.Sprintf(msg, args...)

	f, _ := os.OpenFile(KafkaTraceFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	_, _ = f.Write([]byte(m))
}
