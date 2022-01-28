package main

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/common/clients"
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
	rdc := clients.NewRedisClient(config.RedisAddr, config.RedisPwd, int(config.RedisDB))

	subscribe := rdc.Subscribe(ctx, "new_video")
	channel := subscribe.Channel()

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
