package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func New(cfg *redisConfig, logger *slog.Logger) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	logger.Info("connecting to redis",
		slog.String("addr", addr),
		slog.Int("db", cfg.DB),
	)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("redis connection failed",
			slog.String("addr", addr),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Info("successfully connected to redis", slog.String("addr", addr))

	return client, nil
}
