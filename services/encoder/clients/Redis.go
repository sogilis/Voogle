package clients

import "github.com/go-redis/redis/v8"

func NewRedisClient(addr, pwd string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
}
