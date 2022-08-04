package clients

import (
	"github.com/streadway/amqp"
)

var _ IAmqpClient = amqpClientDummy{}

type amqpClientDummy struct {
	publish      func(string, []byte) error
	consume      func(string) (<-chan amqp.Delivery, error)
	queueDeclare func() (amqp.Queue, error)
}

func NewAmqpClientDummy(publish func(string, []byte) error, consume func(string) (<-chan amqp.Delivery, error), queueDeclare func() (amqp.Queue, error)) IAmqpClient {
	return amqpClientDummy{publish, consume, queueDeclare}
}

func (r amqpClientDummy) Close() error {
	return nil
}

func (r amqpClientDummy) WithRedial() chan IAmqpClient {
	return nil
}

func (r amqpClientDummy) WithExchanger(exchangerName string) error {
	return nil
}

func (r amqpClientDummy) QueueBind(nameQueue string, routingKey string) error {
	return nil
}

func (r amqpClientDummy) GetRandomQueueName() string {
	return ""
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

func (r amqpClientDummy) QueueDeclare() (amqp.Queue, error) {
	if r.queueDeclare != nil {
		return r.queueDeclare()
	}
	queue := amqp.Queue{}

	return queue, nil
}
