package eventhandler

import (
	"context"
	"database/sql"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
)

func ConsumeEvents(amqpClientVideoEncode clients.IAmqpClient, db *sql.DB) {

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
				log.Error("Fail to unmarshal video event")
				continue
			}

			log.Debug("New message received: ", videoProto)
			video := protobuf.VideoProtobufToVideo(videoProto)

			// Update videos status : COMPLETE or FAIL_ENCODE
			videoDb, err := dao.GetVideo(context.Background(), db, video.ID)
			if err != nil {
				log.Errorf("Failed to get video %v from database : %v ", video.ID, err)
				continue
			}

			videoDb.Status = video.Status
			if err := dao.UpdateVideo(context.Background(), db, videoDb); err != nil {
				log.Errorf("Unable to update videos with status  %v: %v", videoDb.Status, err)
			}

			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				log.Error("Failed to Ack message ", video.ID, " - ", err)
				continue
			}
		}
	}
}
