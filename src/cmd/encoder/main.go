package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/encoder/config"
	"github.com/Sogilis/Voogle/src/cmd/encoder/encoding"
)

func main() {
	log.Info("Starting Voogle encoder")
	log.SetLevel(log.DebugLevel)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}

	// Redis client that listen to events
	redisClient := clients.NewRedisClient(cfg.RedisAddr, cfg.RedisPwd, int(cfg.RedisDB))
	// Check if we have successfully opened the connection
	if redisClient.Ping(context.Background()) != nil {
		log.Fatal("Failed to create Redis client")
	}
	channel := redisClient.Subscribe(context.Background(), events.VideoUploaded)

	// S3 client to access the videos
	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client ", err)
	}

	for sub := range channel {
		video := &contracts.Video{}

		if err := proto.Unmarshal([]byte(sub.Payload), video); err != nil {
			log.Error("Fail to unmarshal video event")
			continue
		}

		log.Debug("New message received: ", video)
		log.Info("Starting encoding of video with ID ", video.Id)
		if err := encoding.Process(s3Client, video); err != nil {
			log.Error("Failed to processing video ", video.Id, " - ", err)
		}
		log.Info("Successfully encoded video with ID ", video.Id)

	}
}
