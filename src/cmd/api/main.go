package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Crashed with error: ", err)
		}
	}()

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
