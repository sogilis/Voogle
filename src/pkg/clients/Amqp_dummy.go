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

//-------------------------------------------------

var _ IAmqpExchanger = amqpExchangerDummy{}

type amqpExchangerDummy struct {
	publish      func(string, []byte) error
	consume      func(amqp.Queue) (<-chan amqp.Delivery, error)
	queueDeclare func() (amqp.Queue, error)
	queueBind    func(amqp.Queue, string) error
}

func NewAmqpExchangeDummy(publish func(string, []byte) error, consume func(amqp.Queue) (<-chan amqp.Delivery, error), queueDeclare func() (amqp.Queue, error), queueBind func(amqp.Queue, string) error) IAmqpExchanger {
	return amqpExchangerDummy{publish, consume, queueDeclare, queueBind}
}

func (r amqpExchangerDummy) Publish(nameQueue string, message []byte) error {
	if r.publish != nil {
		return r.publish(nameQueue, message)
	}
	return nil
}

func (r amqpExchangerDummy) Consume(q amqp.Queue) (<-chan amqp.Delivery, error) {
	if r.consume != nil {
		return r.consume(q)
	}
	return nil, nil //nolint:nilnil
}

func (r amqpExchangerDummy) QueueDeclare() (amqp.Queue, error) {
	if r.queueDeclare != nil {
		return r.queueDeclare()
	}
	queue := amqp.Queue{}

	return queue, nil
}

func (r amqpExchangerDummy) QueueBind(q amqp.Queue, key string) error {
	if r.queueBind != nil {
		return r.queueBind(q, key)
	}
	return nil
}
