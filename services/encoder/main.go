package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/services/encoder/clients"
	"github.com/Sogilis/Voogle/services/encoder/config"
	contracts "github.com/Sogilis/Voogle/services/encoder/contracts/v1"
)

var ctx = context.Background()

func main() {
	log.Info("Starting Voogle encoder")

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}

	redisClient := clients.NewRedisClient(config.RedisAddr, config.RedisPwd, int(config.RedisDB))
	err = redisClient.Ping(context.Background())
	if err != nil {
		log.Error("Failed to create Redis client")
	}

	channel := redisClient.Subscribe(ctx, "video_uploaded_on_S3")

	for {
		select {
		case sub := <-channel:
			video := &contracts.Video{}

			err := proto.Unmarshal([]byte(sub.Payload), video)
			if err != nil {
				panic(err)
			}

			fmt.Println(sub)
			fmt.Println(video.GetTitle())
		}
	}
}
