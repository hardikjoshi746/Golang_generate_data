package redis

import (
	"context"

	"github.com/hardikjoshi746/Golang_generate_data/config"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
	cfg := config.Load()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
