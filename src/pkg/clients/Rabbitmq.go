package clients

import (
	"github.com/streadway/amqp"
)

type IRabbitmqClient interface {
	QueueDeclare(name string) (amqp.Queue, error)
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
}

var _ IRabbitmqClient = rabbitmqClient{}

type rabbitmqClient struct {
	rabbitmqClient *amqp.Channel
}

func NewRabbitmqClient(addr, user, pwd string) (IRabbitmqClient, error) {
	amqpConn, err := amqp.Dial("amqp://" + user + ":" + pwd + "@" + addr + "/")
	if err != nil {
		amqpConn.Close()
		return nil, err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		channel.Close()
		return nil, err
	}

	return &rabbitmqClient{
		rabbitmqClient: channel,
	}, nil
}

func (r rabbitmqClient) QueueDeclare(name string) (amqp.Queue, error) {
	return r.rabbitmqClient.QueueDeclare(name, false, false, false, false, nil)
}

func (r rabbitmqClient) Publish(nameQueue string, message []byte) error {
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

func (r rabbitmqClient) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
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
