package di

import (
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultRabbitName = "rabbit"
)

var RabbitDi = &rabbitDi{
	store: make(map[string]*amqp.Channel),
	mutex: sync.RWMutex{},
}

type rabbitDi struct {
	store map[string]*amqp.Channel
	mutex sync.RWMutex
}

func (d *rabbitDi) SetRabbitDi(name string, v *amqp.Channel) {
	d.mutex.Lock()
	d.store[name] = v
	d.mutex.Unlock()
}

func (d *rabbitDi) GetRabbitDi(name string) *amqp.Channel {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.store[name]
}

func GetRabbitWithName(name string) *amqp.Channel {
	return RabbitDi.GetRabbitDi(name)
}

func SetRabbitWithName(name string, v *amqp.Channel) {
	RabbitDi.SetRabbitDi(name, v)
}

func GetRabbit() *amqp.Channel {
	return RabbitDi.GetRabbitDi(DefaultRabbitName)
}

func SetRabbit(v *amqp.Channel) {
	RabbitDi.SetRabbitDi(DefaultRabbitName, v)
}
