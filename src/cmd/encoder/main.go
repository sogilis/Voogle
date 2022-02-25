package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/encoder/config"
	"github.com/Sogilis/Voogle/src/cmd/encoder/encoding"
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

	amqpClient, err := clients.NewAmqpClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoUploaded)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	msgs, err := amqpClient.Consume(events.VideoUploaded)
	if err != nil {
		log.Fatal("Failed to consume RabbitMQ client: ", err)
	}

	for {
		for msg := range msgs {
			video := &contracts.Video{}
			if err := proto.Unmarshal([]byte(msg.Body), video); err != nil {
				log.Error("Fail to unmarshal video event")
				continue
			}

			log.Debug("New message received: ", video)
			log.Info("Starting encoding of video with ID ", video.Id)
			if err := encoding.Process(s3Client, video); err != nil {
				log.Error("Failed to processing video ", video.Id, " - ", err)

				// Nack message but do not requeue it to avoid infinite loop
				if err = msg.Acknowledger.Nack(msg.DeliveryTag, false, false); err != nil {
					log.Error("Failed to Nack message ", video.Id, " - ", err)
					continue
				}
				continue
			}

			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				log.Error("Failed to Ack message ", video.Id, " - ", err)
				continue
			}
		}
		// TODO : if rabbitmq service crash, we get out of the "for msg..." loop
		// and never go back inside.
	}
}
