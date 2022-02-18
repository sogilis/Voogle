package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}
	if cfg.DevMode {
		log.SetLevel(log.DebugLevel)
	}

	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Error("Failed to create S3 client: ", err)
	}

	rabbitmqClient, err := clients.NewRabbitmqClient(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoUploaded)
	if err != nil {
		log.Error("Failed to create RabbitMQ client: ", err)
	}

	routerClients := &router.Clients{
		S3Client:       s3Client,
		RabbitmqClient: rabbitmqClient,
	}

	log.Info("Starting server on port:", cfg.Port)
	srv := &http.Server{
		Handler: router.NewRouter(cfg, routerClients),
		Addr:    fmt.Sprintf("0.0.0.0:%v", cfg.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}

}
