package clients

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type IAmqpClient interface {
	WithRedial() chan IAmqpClient
	Close() error
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
}

var _ IAmqpClient = &amqpClient{}

type amqpClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	user       string
	pwd        string
	address    string
}

func NewAmqpClient(user string, pwd string, addr string) (IAmqpClient, error) {
	amqpC := &amqpClient{
		connection: nil,
		channel:    nil,
		user:       user,
		pwd:        pwd,
		address:    addr,
	}

	conn, err := amqp.Dial("amqp://" + user + ":" + pwd + "@" + addr + "/")
	if err != nil {
		return amqpC, err
	}
	amqpC.connection = conn

	channel, err := conn.Channel()
	if err != nil {
		return amqpC, err
	}
	amqpC.channel = channel

	return amqpC, err
}

// This function return a channel with a new, working client.
// Ensure the previous client is closed so this function can reconnect.
func (r *amqpClient) WithRedial() chan IAmqpClient {
	session := make(chan IAmqpClient)
	go func() {
		// To avoid overcharging the service, we set a 5 seconds timer if an error occur.
		pauseTimer := 5 * time.Second
		defer close(session)
		for {
			client, err := NewAmqpClient(r.user, r.pwd, r.address)
			if err != nil {
				log.Error("Could not reconnect to RabbitMQ : ", err)
				time.Sleep(pauseTimer)
				continue
			}
			session <- client
		}
	}()
	return session
}

func (r *amqpClient) Publish(nameQueue string, message []byte) error {
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

func (r *amqpClient) Consume(nameQueue string) (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(nameQueue, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
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

func (r *amqpClient) Close() error {
	return r.connection.Close()
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
