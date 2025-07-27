package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // e.g. "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // "" if no password
		DB:       0,
	})
	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
