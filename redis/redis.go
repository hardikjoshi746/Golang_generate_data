package redis

import (
	"context"
	"fmt"
	"time"

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

const maxRequestsPerMinute = 10

// RateLimit returns true if the user is rate limited
func RateLimit(userID string) (bool, error) {
	key := fmt.Sprintf("rate:%s", userID)

	// Increment the user's request count
	count, err := Rdb.Incr(Ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Set expiration the first time
	if count == 1 {
		err := Rdb.Expire(Ctx, key, time.Minute).Err()
		if err != nil {
			return false, err
		}
	}

	// Limit exceeded
	if count > maxRequestsPerMinute {
		return true, nil
	}

	return false, nil
}
