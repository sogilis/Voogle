package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"

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

	amqpClientVideoUpload, _ := clients.NewAmqpClient(cfg.RabbitmqUser, cfg.RabbitmqPwd, cfg.RabbitmqAddr)

	// Listen, consume and publish on amqpClientVideoUpload
	eventhandler.ConsumeEvents(amqpClientVideoUpload, s3Client)
}
