package main

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	. "github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	config, err := NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}
	if config.IsDev {
		log.SetLevel(log.DebugLevel)
	}

	s3Client, err := clients.NewS3Client(config.S3Host, config.S3Region, config.S3Bucket, config.S3AuthKey, config.S3AuthPwd)
	if err != nil {
		log.Error("Failed to create S3 client")
	}

	redisClient := clients.NewRedisClient(config.RedisAddr, config.RedisPwd, config.RedisDB)
	err = redisClient.Ping(context.Background())
	if err != nil {
		log.Error("Failed to create Redis client")
	}

	routerClients := &router.Clients{
		S3Client:    s3Client,
		RedisClient: redisClient,
	}

	log.Info("Starting server on port:", config.Port)
	srv := &http.Server{
		Handler: router.NewRouter(config, routerClients),
		Addr:    fmt.Sprintf("0.0.0.0:%v", config.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}

}
