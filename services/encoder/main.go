package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

type Video struct {
	Title string `json:"title"`
}

func main() {
	log.Info("Starting Voogle encoder")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	subscribe := rdb.Subscribe(ctx, "new_video")
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
