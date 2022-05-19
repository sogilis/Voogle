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
		// Only consumer should declare queue
		// if _, err := r.QueueDeclare(); err != nil {
		// 	return err
		// }

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

type IAmqpExchanger interface {
	Publish(key string, message []byte) error
	Consume(q amqp.Queue) (<-chan amqp.Delivery, error)
	QueueDeclare() (amqp.Queue, error)
	QueueBind(q amqp.Queue, key string) error
}

var _ IAmqpExchanger = &amqpExchanger{}

type amqpExchanger struct {
	channel       *amqp.Channel
	fullAddr      string
	exchangerName string
}

func NewAmqpExchanger(addr, user, pwd, exchangerName string) (IAmqpExchanger, error) {
	amqpExchanger := &amqpExchanger{
		channel:       nil,
		fullAddr:      "amqp://" + user + ":" + pwd + "@" + addr + "/",
		exchangerName: exchangerName,
	}

	if err := amqpExchanger.connect(); err != nil {
		return nil, err
	}

	if err := amqpExchanger.channel.ExchangeDeclare(
		amqpExchanger.exchangerName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return amqpExchanger, nil
}

func (r *amqpExchanger) connect() error {
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

func (r *amqpExchanger) Publish(key string, message []byte) error {
	err := r.channel.Publish(
		r.exchangerName,
		key,
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
		// Only consumer should declare queue
		// if _, err := r.QueueDeclare(); err != nil {
		// 	return err
		// }

		return r.channel.Publish(
			r.exchangerName,
			key,
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

func (r *amqpExchanger) Consume(q amqp.Queue) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *amqpExchanger) QueueDeclare() (amqp.Queue, error) {
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	return q, err
}

func (r *amqpExchanger) QueueBind(q amqp.Queue, key string) error {
	err := r.channel.QueueBind(
		q.Name,
		key,
		r.exchangerName,
		false,
		nil)
	return err
}
