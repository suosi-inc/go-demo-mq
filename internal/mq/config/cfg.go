package config

// Cfg is singleton
var Cfg = &cfg{}

type cfg struct {
	App         app         `json:"app"`
	Logger      logger      `json:"logger"`
	Rabbit      rabbit      `json:"rabbit"`
	RabbitQueue rabbitQueue `json:"rabbitQueue"`
	Kafka       kafka       `json:"kafka"`
	KafkaQueue  kafkaQueue  `json:"kafkaQueue"`
}

type app struct {
	Name string `json:"name"`
}

type logger struct {
	File       string `json:"file"`
	Format     string `json:"format"`
	Level      string `json:"level"`
	MaxSize    int    `json:"maxSize"`
	MaxAge     int    `json:"maxAge"`
	MaxBackups int    `json:"maxBackups"`
}

type rabbit struct {
	Amqp         string `json:"amqp"`
	Prefetch     int    `json:"prefetch"`
	DeliveryMode uint8  `json:"deliveryMode"`
}

type rabbitQueue struct {
	Simple struct {
		Name       string `json:"name"`
		Goroutines int    `json:"goroutines"`
	} `json:"simple"`
	Topic struct {
		Exchange    string   `json:"exchange"`
		Name        string   `json:"name"`
		Goroutines  int      `json:"goroutines"`
		RoutingKeys []string `json:"routingKeys"`
	}
}

type kafka struct {
	Servers  []string `json:"servers"`
	Compress string   `json:"compress"`
	Sasl     struct {
		Enable    bool   `json:"enable"`
		Mechanism string `json:"mechanism"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	} `json:"sasl"`
}

type kafkaQueue struct {
	Topic      string `json:"topic"`
	Goroutines int    `json:"goroutines"`
	GroupId    string `json:"groupId"`
}
