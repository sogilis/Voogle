package clients

import (
	"github.com/Sogilis/Voogle/src/pkg/events"
	"github.com/streadway/amqp"
)

type IRabbitmqClient interface {
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
}

var _ IRabbitmqClient = &rabbitmqClient{}

type rabbitmqClient struct {
	channel  *amqp.Channel
	fullAddr string
}

func NewRabbitmqClient(addr, user, pwd, queueName string) (IRabbitmqClient, error) {
	rabbitmqC := &rabbitmqClient{
		channel:  nil,
		fullAddr: "amqp://" + user + ":" + pwd + "@" + addr + "/",
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

func (r *rabbitmqClient) queueDeclare(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(name, false, false, false, false, nil)
}

func (r *rabbitmqClient) Publish(nameQueue string, message []byte) error {
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

func (r *rabbitmqClient) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
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
