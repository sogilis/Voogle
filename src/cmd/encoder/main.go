package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/encoder/config"
	"github.com/Sogilis/Voogle/src/cmd/encoder/eventhandler"
)

func main() {
	log.Info("Starting Voogle encoder")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}
	if cfg.DevMode {
		log.SetLevel(log.DebugLevel)
	}

	// S3 client to access the videos
	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client ", err)
	}

	// amqpClient for new uploaded video (api->encoder)
	amqpClientVideoUpload, err := clients.NewAmqpClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoUploaded)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	// Consumer only should declare queue
	if _, err := amqpClientVideoUpload.QueueDeclare(); err != nil {
		log.Fatal("Failed to declare RabbitMQ queue: ", err)
	}

	// amqpClient for encoded video (encoder->api)
	amqpClientVideoEncode, err := clients.NewAmqpClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoEncoded)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	// Listen and consume on amqpClientVideoUpload and publish video status on amqpClientVideoEncode
	eventhandler.ConsumeEvents(amqpClientVideoUpload, s3Client, amqpClientVideoEncode)
}
