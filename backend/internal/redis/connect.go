package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Init() *redis.Client {
	rbd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := rbd.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return rbd
}
