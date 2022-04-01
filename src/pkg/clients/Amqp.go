package clients

import (
	"github.com/streadway/amqp"
)

type IAmqpClient interface {
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
	QueueDeclare() (amqp.Queue, error)
}

var _ IAmqpClient = &amqpClient{}

type amqpClient struct {
	channel   *amqp.Channel
	fullAddr  string
	queueName string
}

func NewAmqpClient(addr, user, pwd, queueName string) (IAmqpClient, error) {
	amqpC := &amqpClient{
		channel:   nil,
		fullAddr:  "amqp://" + user + ":" + pwd + "@" + addr + "/",
		queueName: queueName,
	}

	if err := amqpC.connect(); err != nil {
		return nil, err
	}

	return amqpC, nil
}

func (r *amqpClient) Publish(nameQueue string, message []byte) error {
	err := r.channel.Publish(
		"",
		nameQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)

	if err != nil {
		if err = r.connect(); err != nil {
			return err
		}

		// If we cannot publish, try to reconnect to rabbitMQ service ONE time before
		// return error
		if _, err := r.QueueDeclare(); err != nil {
			return err
		}

		return r.channel.Publish(
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
	return nil
}

func (r *amqpClient) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		nameQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *amqpClient) connect() error {
	amqpConn, err := amqp.Dial(r.fullAddr)
	if err != nil {
		return err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return err
	}

	r.channel = channel
	return nil
}

func (r *amqpClient) QueueDeclare() (amqp.Queue, error) {
	return r.channel.QueueDeclare(r.queueName, false, false, false, false, nil)
}
