package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}
	if cfg.DevMode {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}

	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Error("Failed to create S3 client: ", err)
	}

	// amqpClient for new uploaded video (api->encoder)
	amqpClientVideoUpload, err := clients.NewAmqpClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoUploaded)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	// amqpClient for encoded video (encoder->api)
	amqpClientVideoEncode, err := clients.NewAmqpClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoEncoded)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	// Consumer only should declare queue
	if _, err := amqpClientVideoEncode.QueueDeclare(); err != nil {
		log.Fatal("Failed to declare RabbitMQ queue: ", err)
	}

	// Use "?parseTime=true" to match golang time.Time with Mariadb DATETIME types
	db, err := sql.Open("mysql", cfg.MariadbUser+":"+cfg.MariadbUserPwd+"@tcp("+cfg.MariadbAddr+")/"+cfg.MariadbName+"?parseTime=true")
	if err != nil {
		log.Fatal("Failed to open connection with database: ", err)
	}
	defer db.Close()

	routerClients := &router.Clients{
		S3Client:      s3Client,
		AmqpClient:    amqpClientVideoUpload,
		MariadbClient: db,
	}

	uuidGen := uuidgenerator.NewUuidGenerator()

	routerUUIDGen := &router.UUIDGenerator{
		UUIDGen: uuidGen,
	}

	log.Info("Starting server on port:", cfg.Port)
	srv := &http.Server{
		Handler: router.NewRouter(cfg, routerClients, routerUUIDGen),
		Addr:    fmt.Sprintf("0.0.0.0:%v", cfg.Port),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Crashed with error: ", err)
		}
	}()

	// Listen to encoder video status update
	msgs, err := amqpClientVideoEncode.Consume(events.VideoEncoded)
	if err != nil {
		log.Fatal("Failed to consume RabbitMQ client: ", err)
	}

	go func() { consumeEvents(msgs, db) }()

	c := make(chan os.Signal, 1)

	// Catch SIGINT, SIGKILL, SIGQUIT or SIGTERM
	signal.Notify(c, os.Interrupt)

	sig := waitInterruptSignal(c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Info("HTTP server Shutdown: ", err)
	}

	log.Infof("Receive signal %v. Shutting down properly", sig)
}

func waitInterruptSignal(ch <-chan os.Signal) os.Signal {
	return <-ch
}

func consumeEvents(msgs <-chan amqp.Delivery, db *sql.DB) {
	for {
		for msg := range msgs {
			video := &contracts.EncodedVideo{}
			if err := proto.Unmarshal([]byte(msg.Body), video); err != nil {
				log.Error("Fail to unmarshal video event")
				continue
			}

			log.Debug("New message received: ", video)

			// Update videos status : COMPLETE or FAIL_ENCODE
			videoDb, err := dao.GetVideo(db, video.Id)
			if err != nil {
				log.Errorf("Failed to get video %v from database : %v ", video.Id, err)
				continue
			}
			videoDb.Status = models.VideoStatus(video.Status)
			log.Debug("Update video")
			if err := dao.UpdateVideo(db, videoDb); err != nil {
				log.Errorf("Unable to update videos with status  %v: %v", videoDb.Status, err)
			}

			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				log.Error("Failed to Ack message ", video.Id, " - ", err)
				continue
			}
		}
	}
}
