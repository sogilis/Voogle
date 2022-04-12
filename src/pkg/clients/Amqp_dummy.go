package clients

import (
	"github.com/streadway/amqp"
)

var _ IAmqpClient = amqpClientDummy{}

type amqpClientDummy struct {
	publish func(string, []byte) error
	consume func(string) (<-chan amqp.Delivery, error)
}

func NewAmqpClientDummy(publish func(string, []byte) error, consume func(string) (<-chan amqp.Delivery, error)) IAmqpClient {
	return amqpClientDummy{publish, consume}
}

func (r amqpClientDummy) Publish(nameQueue string, message []byte) error {
	if r.publish != nil {
		return r.publish(nameQueue, message)
	}
	return nil
}

func (r amqpClientDummy) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
	if r.consume != nil {
		return r.consume(nameQueue)
	}
	return nil, nil //nolint:nilnil
}
