package eventhandler

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

func ConsumeEvents(cfg config.Config, amqpExchangerStatus clients.IAmqpExchanger, videosDAO *dao.VideosDAO) {
	// amqpClient for encoded video (encoder->api)
	amqpClientVideoEncode, err := clients.NewAmqpClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoEncoded)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	// Consumer only should declare queue
	if _, err := amqpClientVideoEncode.QueueDeclare(); err != nil {
		log.Fatal("Failed to declare RabbitMQ queue: ", err)
	}

	// Listen to encoder video status update
	msgs, err := amqpClientVideoEncode.Consume(events.VideoEncoded)
	if err != nil {
		log.Fatal("Failed to consume RabbitMQ client: ", err)
	}

	for {
		for msg := range msgs {
			videoProto := &contracts.Video{}
			if err := proto.Unmarshal([]byte(msg.Body), videoProto); err != nil {
				log.Error("Fail to unmarshal video event : ", err)
				continue
			}

			log.Debug("New message received: ", videoProto)
			video := protobuf.VideoProtobufToVideo(videoProto)

			// Update videos status : COMPLETE or FAIL_ENCODE
			videoDb, err := videosDAO.GetVideo(context.Background(), video.ID)
			if err != nil {
				log.Errorf("Failed to get video %v from database : %v ", video.ID, err)
				continue
			}

			videoDb.Status = video.Status
			videoDb.CoverPath = video.CoverPath
			if err := videosDAO.UpdateVideo(context.Background(), videoDb); err != nil {
				log.Errorf("Unable to update videos with status  %v: %v", videoDb.Status, err)
			}
			if video.Status == models.COMPLETE {
				metrics.CounterVideoEncodeSuccess.Inc()
			} else if video.Status == models.FAIL_ENCODE {
				metrics.CounterVideoEncodeFail.Inc()
			}

			publishStatus(amqpExchangerStatus, videoDb)

			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				log.Error("Failed to Ack message ", video.ID, " - ", err)
				continue
			}
		}
	}
}

func publishStatus(amqpExchanger clients.IAmqpExchanger, video *models.Video) {
	msg, err := proto.Marshal(protobuf.VideoToVideoProtobuf(video))
	if err != nil {
		log.Error("Failed to Marshal status", err)
	}
	if err := amqpExchanger.Publish(video.Title, msg); err != nil {
		log.Error("Unable to publish status update", err)
	}
}
