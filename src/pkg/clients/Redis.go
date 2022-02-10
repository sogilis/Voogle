package clients

import (
	"context"
	"crypto/tls"

	"github.com/go-redis/redis/v8"
)

type IRedisClient interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) <-chan *redis.Message
	Ping(ctx context.Context) error
}

var _ IRedisClient = redisClient{}

type redisClient struct {
	redisClient *redis.ClusterClient
}

func NewRedisClient(addr, pwd string, useTLS bool) IRedisClient {

	tlsConfig := &tls.Config{}
	if !useTLS {
		tlsConfig = nil
	}

	return &redisClient{
		redisClient: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     []string{addr},
			Password:  pwd,
			TLSConfig: tlsConfig,
		}),
	}
}

func (r redisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.redisClient.Publish(ctx, channel, message).Err()
}

func (r redisClient) Subscribe(ctx context.Context, channel string) <-chan *redis.Message {
	subscribe := r.redisClient.Subscribe(ctx, channel)
	return subscribe.Channel()
}

func (r redisClient) Ping(ctx context.Context) error {
	_, err := r.redisClient.Ping(ctx).Result()
	return err
}
