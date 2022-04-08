package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
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

	msgs, err := amqpClientVideoUpload.Consume(events.VideoUploaded)
	if err != nil {
		log.Fatal("Failed to consume RabbitMQ client: ", err)
	}

	// Nack message but do not requeue it to avoid infinite loop
	// TODO : if rabbitmq service crash, we get out of the "for msg..." loop
	// and never go back inside.
	consumeEvents(msgs, s3Client, amqpClientVideoEncode)
}

func consumeEvents(msgs <-chan amqp.Delivery, s3Client clients.IS3Client, amqpClientVideoEncode clients.IAmqpClient) {
	for {
		for msg := range msgs {
			video := &contracts.UploadedVideo{}
			if err := proto.Unmarshal([]byte(msg.Body), video); err != nil {
				log.Error("Fail to unmarshal video event")
				continue
			}

			videoEncoded := &contracts.EncodedVideo{
				Id:     video.Id,
				Status: int32(models.ENCODING),
			}

			log.Debug("New message received: ", video)
			log.Info("Starting encoding of video with ID ", video.Id)

			if err := encoding.Process(s3Client, video); err != nil {
				log.Error("Failed to processing video ", video.Id, " - ", err)

				if err = msg.Acknowledger.Nack(msg.DeliveryTag, false, false); err != nil {
					log.Error("Failed to Nack message ", video.Id, " - ", err)
				}

				// Send video status updated : FAIL_ENCODE
				videoEncoded.Status = int32(models.FAIL_ENCODE)
				if err = sendUpdatedVideoStatus(videoEncoded, amqpClientVideoEncode); err != nil {
					log.Error("Error while sending new video status : ", err)
				}

				continue
			}

			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				log.Error("Failed to Ack message ", video.Id, " - ", err)

				// Send video status updated : FAIL_ENCODE
				videoEncoded.Status = int32(models.FAIL_ENCODE)
				if err = sendUpdatedVideoStatus(videoEncoded, amqpClientVideoEncode); err != nil {
					log.Error("Error while sending new video status : ", err)
				}

				continue
			}

			// Send video status updated : COMPLETE
			videoEncoded.Status = int32(models.COMPLETE)
			if err := sendUpdatedVideoStatus(videoEncoded, amqpClientVideoEncode); err != nil {
				log.Error("Error while sending new video status : ", err)
				continue
			}
		}
	}
}

func sendUpdatedVideoStatus(video *contracts.EncodedVideo, amqpC clients.IAmqpClient) error {
	videoData, err := proto.Marshal(video)
	if err != nil {
		log.Error("Unable to marshal video ", err)
		return err
	}

	if err = amqpC.Publish(events.VideoEncoded, videoData); err != nil {
		log.Error("Unable to publish on Amqp client VideoEncode ", err)
		return err
	}

	return nil
}
