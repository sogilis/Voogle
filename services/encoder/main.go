package main

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/encoder/clients"
	. "github.com/Sogilis/Voogle/services/encoder/config"
)

var ctx = context.Background()

type Video struct {
	Title string `json:"title"`
}

func main() {
	log.Info("Starting Voogle encoder")

	config, err := NewConfig()
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
			video := &Video{}

			err := json.Unmarshal([]byte(sub.Payload), video)
			if err != nil {
				panic(err)
			}

			fmt.Println(sub)
			fmt.Println(video.Title)
		}
	}
}
