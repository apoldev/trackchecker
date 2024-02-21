package redis

import (
	"fmt"

	"github.com/apoldev/trackchecker/internal/app/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg *config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   0,
	})
}
