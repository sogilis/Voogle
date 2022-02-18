package clients

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type IRabbitmqClient interface {
	QueueDeclare(name string) (amqp.Queue, error)
	Publish(nameQueue string, message []byte) error
	Consume(nameQueue string) (<-chan amqp.Delivery, error)
	Connect() error
	Reconnect()
}

var _ IRabbitmqClient = &rabbitmqClient{}
var fullAddr string
var queueName string

type rabbitmqClient struct {
	rabbitmqClient *amqp.Channel
}

func NewRabbitmqClient(addr, user, pwd string) (IRabbitmqClient, error) {
	fullAddr = "amqp://" + user + ":" + pwd + "@" + addr + "/"
	rabbitmqC := &rabbitmqClient{
		rabbitmqClient: nil,
	}
	if err := rabbitmqC.Connect(); err != nil {
		return nil, err
	}
	return rabbitmqC, nil
}

func (r *rabbitmqClient) Connect() error {
	amqpConn, err := amqp.Dial(fullAddr)
	if err != nil {
		//amqpConn.Close()
		return err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		//channel.Close()
		return err
	}

	r.rabbitmqClient = channel
	return nil
}

func (r *rabbitmqClient) QueueDeclare(name string) (amqp.Queue, error) {
	queueName = name
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

func (r *rabbitmqClient) Reconnect() {
	go func() {
		// time.Sleep(10 * time.Second)
		log.Infof("Closing: %s", <-r.rabbitmqClient.NotifyClose(make(chan *amqp.Error)))
		for err := r.Connect(); err != nil; err = r.Connect() {
			log.Info(err)
			time.Sleep(10 * time.Second)
		}
		log.Info("Reconnected")
		r.QueueDeclare(queueName)
	}()
}
