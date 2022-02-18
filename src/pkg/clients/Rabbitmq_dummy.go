package clients

import (
	"github.com/streadway/amqp"
)

var _ IRabbitmqClient = rabbitmqClientDummy{}

type rabbitmqClientDummy struct {
	publish func(string, []byte) error
	consume func(string) (<-chan amqp.Delivery, error)
}

func NewRabbitmqClientDummy(publish func(string, []byte) error, consume func(string) (<-chan amqp.Delivery, error)) IRabbitmqClient {
	return rabbitmqClientDummy{publish, consume}
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
