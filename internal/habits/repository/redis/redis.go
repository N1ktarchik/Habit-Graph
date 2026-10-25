package redis

import (
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	log   *slog.Logger
	redis *redis.Client
}

func NewCache(log *slog.Logger, redis *redis.Client) *cache {
	return &cache{
		log:   log,
		redis: redis,
	}
}

func buildKey(id, year string) string {
	if year == "" {
		return fmt.Sprintf("habit:%s:graph:*", id)
	}

	return fmt.Sprintf("habit:%s:graph:%s", id, year)

}
