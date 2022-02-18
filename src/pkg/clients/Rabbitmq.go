package clients

import (
	"github.com/streadway/amqp"
)

type IRabbitmqClient interface {
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
}

var _ IRabbitmqClient = &rabbitmqClient{}
var fullAddr string

type rabbitmqClient struct {
	rabbitmqClient *amqp.Channel
}

func NewRabbitmqClient(addr, user, pwd, queueName string) (IRabbitmqClient, error) {
	fullAddr = "amqp://" + user + ":" + pwd + "@" + addr + "/"
	rabbitmqC := &rabbitmqClient{
		rabbitmqClient: nil,
	}

	if err := rabbitmqC.connect(); err != nil {
		return nil, err
	}

	if _, err := rabbitmqC.queueDeclare(queueName); err != nil {
		return nil, err
	}

	return rabbitmqC, nil
}

func (r *rabbitmqClient) connect() error {
	amqpConn, err := amqp.Dial(fullAddr)
	if err != nil {
		return err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return err
	}

	r.rabbitmqClient = channel
	return nil
}

func (r *rabbitmqClient) queueDeclare(name string) (amqp.Queue, error) {
	return r.rabbitmqClient.QueueDeclare(name, false, false, false, false, nil)
}

func (r *rabbitmqClient) Publish(nameQueue string, message []byte) error {
	return r.rabbitmqClient.Publish(
		"",
		nameQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

func (r *rabbitmqClient) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
	return r.rabbitmqClient.Consume(
		nameQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
