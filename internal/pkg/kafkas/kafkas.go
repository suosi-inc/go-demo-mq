package kafkas

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/suosi-inc/go-demo/mq/internal/mq/config"
	"github.com/x-funs/go-fun"
)

func NewWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{}

	servers := config.Cfg.Kafka.Servers

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
			Addr:      kafka.TCP(servers...),
			Topic:     topic,
			Balancer:  &kafka.Hash{},
			Transport: sharedTransport,
		}

	} else {
		w = &kafka.Writer{
			Addr:        kafka.TCP(servers...),
			Topic:       topic,
			Compression: compress,
			Balancer:    &kafka.Hash{},
		}
	}

	return w
}

func NewReader(topic string, groupId string) *kafka.Reader {
	r := &kafka.Reader{}

	servers := config.Cfg.Kafka.Servers

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

		r = kafka.NewReader(kafka.ReaderConfig{
			Brokers: servers,
			GroupID: "consumer-group-id",
			Topic:   "topic-A",
			Dialer:  dialer,
		})

	} else {
		r = kafka.NewReader(kafka.ReaderConfig{
			Brokers: servers,
			GroupID: groupId,
			Topic:   topic,
		})
	}

	return r
}
