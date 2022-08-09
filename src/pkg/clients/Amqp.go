package clients

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AmqpClient interface {
	WithRedial() chan AmqpClient
	WithExchanger(exchangerName string) error
	Close() error
	Publish(routingKey string, message []byte) error
	GetRandomQueueName() string
	QueueBind(nameQueue string, routingKey string) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
}

var _ AmqpClient = &amqpClient{}

type amqpClient struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	user          string
	pwd           string
	address       string
	exchangerName string
}

func NewAmqpClient(user string, pwd string, addr string) (AmqpClient, error) {
	amqpC := &amqpClient{
		connection:    nil,
		channel:       nil,
		user:          user,
		pwd:           pwd,
		address:       addr,
		exchangerName: "",
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
func (r *amqpClient) WithRedial() chan AmqpClient {
	session := make(chan AmqpClient)
	go func() {
		defer close(session)
		for {
			client, err := NewAmqpClient(r.user, r.pwd, r.address)
			if err != nil {
				log.Error("Could not reconnect to RabbitMQ : ", err)
				continue
			}
			if r.exchangerName != "" {
				err := client.WithExchanger(r.exchangerName)
				if err != nil {
					log.Info("Could not create exchanger : ", err)
					continue
				}
			}
			session <- client
		}
	}()
	return session
}

func (r *amqpClient) WithExchanger(exchangerName string) error {
	r.exchangerName = exchangerName
	err := r.channel.ExchangeDeclare(
		exchangerName, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	return err
}

func (r *amqpClient) GetRandomQueueName() string {
	hostname, err := os.Hostname()
	h := sha1.New()
	fmt.Fprint(h, hostname)
	fmt.Fprint(h, err)
	fmt.Fprint(h, os.Getpid())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r *amqpClient) Publish(routingKey string, message []byte) error {
	err := r.channel.Publish(
		r.exchangerName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		conn, err := amqp.Dial(r.address)
		if err != nil {
			return err
		}
		r.connection = conn

		channel, err := conn.Channel()
		if err != nil {
			return err
		}
		r.channel = channel

		// If we cannot publish, try to reconnect to rabbitMQ service ONE time before
		// return error
		// Only consumer should declare queue
		// if _, err := r.QueueDeclare(); err != nil {
		// 	return err
		// }

		return r.channel.Publish(
			r.exchangerName,
			routingKey,
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

func (r *amqpClient) QueueBind(nameQueue string, routingKey string) error {
	if r.exchangerName == "" {
		return errors.New("No exchanger set on this client.")
	}
	err := r.channel.QueueBind(
		nameQueue,
		routingKey,
		r.exchangerName,
		false,
		nil)
	return err
}

func (r *amqpClient) Close() error {
	return r.connection.Close()
}
