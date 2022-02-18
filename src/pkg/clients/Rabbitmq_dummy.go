package clients

import (
	"github.com/streadway/amqp"
)

var _ IRabbitmqClient = rabbitmqClientDummy{}

type rabbitmqClientDummy struct {
	connect      func() error
	reconnect    func()
	queueDeclare func(string) (amqp.Queue, error)
	publish      func(string, []byte) error
	consume      func(string) (<-chan amqp.Delivery, error)
}

func NewRabbitmqClientDummy(connect func() error, reconnect func(), queueDeclare func(string) (amqp.Queue, error), publish func(string, []byte) error, consume func(string) (<-chan amqp.Delivery, error)) IRabbitmqClient {
	return rabbitmqClientDummy{connect, reconnect, queueDeclare, publish, consume}
}

func (r rabbitmqClientDummy) Connect() error {
	if r.connect != nil {
		return r.connect()
	}
	return nil
}

func (r rabbitmqClientDummy) Reconnect() {
	if r.reconnect != nil {
		r.reconnect()
	}
}

func (r rabbitmqClientDummy) QueueDeclare(name string) (amqp.Queue, error) {
	if r.queueDeclare != nil {
		return r.queueDeclare(name)
	}
	return amqp.Queue{}, nil
}

func (r rabbitmqClientDummy) Publish(nameQueue string, message []byte) error {
	if r.publish != nil {
		return r.publish(nameQueue, message)
	}
	return nil
}

func (r rabbitmqClientDummy) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
	if r.consume != nil {
		return r.consume(nameQueue)
	}
	return nil, nil
}
