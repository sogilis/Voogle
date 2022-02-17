package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	. "github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	cfg, err := NewConfig()
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

	redisClient := clients.NewRedisClient(cfg.RedisAddr, cfg.RedisPwd, !cfg.DevMode)
	if err := redisClient.Ping(context.Background()); err != nil {
		log.Error("Failed to create Redis client: ", err)
	}

	routerClients := &router.Clients{
		S3Client:    s3Client,
		RedisClient: redisClient,
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

	// Setup wait time for gracefull shutdown
	wait := time.Second * 15
	c := make(chan os.Signal, 1)

	// Catch SIGINT, SIGKILL, SIGQUIT or SIGTERM
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Failed to shutdown ", err)
	}

	log.Println("shutting down")
	os.Exit(0)
}
