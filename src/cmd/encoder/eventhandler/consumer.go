package eventhandler

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/encoder/encoding"
)

func ConsumeEvents(amqpClientVideoUpload clients.IAmqpClient, s3Client clients.IS3Client) {
	session := amqpClientVideoUpload.WithRedial()
	for {
		client := <-session

		msgs, err := client.Consume(events.VideoUploaded)
		if err != nil {
			log.Error("Failed to consume RabbitMQ client: ", err)
			continue
		}

		for msg := range msgs {
			video := &contracts.Video{}
			if err := proto.Unmarshal([]byte(msg.Body), video); err != nil {
				log.Error("Fail to unmarshal video event : ", err)
				continue
			}

			videoEncoded := &contracts.Video{
				Id:        video.Id,
				Status:    contracts.Video_VIDEO_STATUS_ENCODING,
				Source:    video.Source,
				CoverPath: video.CoverPath,
			}

			log.Debug("New message received: ", video)
			log.Info("Starting encoding of video with ID ", video.Id)

			if err := encoding.Process(s3Client, video); err != nil {
				log.Error("Failed to processing video ", video.Id, " - ", err)

				if err = msg.Acknowledger.Nack(msg.DeliveryTag, false, false); err != nil {
					log.Error("Failed to Nack message ", video.Id, " - ", err)
				}

				// Send video status updated : FAIL_ENCODE
				videoEncoded.Status = contracts.Video_VIDEO_STATUS_FAIL_ENCODE
				if err = sendUpdatedVideoStatus(videoEncoded, client); err != nil {
					log.Error("Error while sending new video status : ", err)
				}

				continue
			}

			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				log.Error("Failed to Ack message ", video.Id, " - ", err)

				// Send video status updated : FAIL_ENCODE
				videoEncoded.Status = contracts.Video_VIDEO_STATUS_FAIL_ENCODE
				if err = sendUpdatedVideoStatus(videoEncoded, client); err != nil {
					log.Error("Error while sending new video status : ", err)
				}

				continue
			}

			// Send updates
			// Update video status to COMPLETE
			videoEncoded.Status = contracts.Video_VIDEO_STATUS_COMPLETE
			// Update video cover path
			if len(videoEncoded.CoverPath) > 0 && filepath.Ext(videoEncoded.CoverPath) != ".jpeg" {
				videoEncoded.CoverPath = videoEncoded.Id + "/cover.jpeg"
			}
			if err := sendUpdatedVideoStatus(videoEncoded, client); err != nil {
				log.Error("Error while sending new video status : ", err)
				continue
			}
		}
		// We close the client to let another take his place.
		client.Close()
	}
}

func sendUpdatedVideoStatus(video *contracts.Video, amqpC clients.IAmqpClient) error {
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
