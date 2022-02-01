package clients

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var _ IRedisClient = redisClientDummy{}

type redisClientDummy struct {
	publish   func(string, interface{}) error
	subscribe func() <-chan *redis.Message
	ping      func() error
}

func NewRedisClientDummy(publish func(string, interface{}) error, subscribe func() <-chan *redis.Message, ping func() error) IRedisClient {
	return redisClientDummy{publish, subscribe, ping}
}

func (r redisClientDummy) Publish(ctx context.Context, channel string, message interface{}) error {
	if r.publish != nil {
		return r.publish(channel, message)
	}
	return nil
}

func (r redisClientDummy) Subscribe(ctx context.Context, channel string) <-chan *redis.Message {
	if r.subscribe != nil {
		return r.subscribe()
	}
	return nil
}

func (r redisClientDummy) Ping(ctx context.Context) error {
	if r.ping != nil {
		return r.ping()
	}
	return nil
}
