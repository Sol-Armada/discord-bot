package store

import (
	"context"

	"github.com/go-redis/redis/v9"
)

var Client *redis.Client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func Exists(ctx context.Context, key string) bool {
	r, err := Client.Get(ctx, key).Result()
	return err == nil && r != ""
}
