package config

// Cfg is singleton
var Cfg = &cfg{}

type cfg struct {
	App    app    `json:"app"`
	Logger logger `json:"logger"`
	Rabbit rabbit `json:"rabbit"`
	Queue  queue  `json:"queue"`
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
	Amqp     string `json:"amqp"`
	Prefetch int    `json:"prefetch"`
}

type queue struct {
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
