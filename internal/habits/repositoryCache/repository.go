package repositoryCache

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

const habitKey = "habit:coding:2026"

type repository struct {
	log   *slog.Logger
	redis *redis.Client
}

func NewRepository(log *slog.Logger, redis *redis.Client) *repository {
	return &repository{
		log:   log,
		redis: redis,
	}
}
