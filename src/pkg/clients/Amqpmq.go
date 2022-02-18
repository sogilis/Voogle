package clients

import (
	"github.com/streadway/amqp"

	"github.com/Sogilis/Voogle/src/pkg/events"
)

type IAmqpClient interface {
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
}

var _ IAmqpClient = &amqpClient{}

type amqpClient struct {
	channel  *amqp.Channel
	fullAddr string
}

func NewAmqpClient(addr, user, pwd, queueName string) (IAmqpClient, error) {
	amqpC := &amqpClient{
		channel:  nil,
		fullAddr: "amqp://" + user + ":" + pwd + "@" + addr + "/",
	}

	if err := amqpC.connect(); err != nil {
		return nil, err
	}

	if _, err := amqpC.queueDeclare(queueName); err != nil {
		return nil, err
	}

	return amqpC, nil
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

func (r *amqpClient) queueDeclare(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(name, false, false, false, false, nil)
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

	// If we cannot publish, try to reconnect to rabbitMQ service ONE time before
	// return error
	if err != nil {
		if err = r.connect(); err != nil {
			return err
		}

		if _, err := r.queueDeclare(events.VideoUploaded); err != nil {
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
		true,
		false,
		false,
		false,
		nil,
	)
}
