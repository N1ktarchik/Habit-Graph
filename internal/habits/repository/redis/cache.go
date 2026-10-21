package redis

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
)

func (c *cache) GetGraph(ctx context.Context, habitID, year string) (map[string]int, error) {
	key := buildKey(habitID, year)

	dataStr, err := c.redis.HGetAll(ctx, key).Result()

	if err != nil {
		c.log.Error("error in repository : redis error : get graph error", slog.Any("err:", err))
		return nil, err
	}

	if len(dataStr) == 0 {
		return nil, nil
	}

	graph := make(map[string]int, len(dataStr))
	for dateKey, minStr := range dataStr {
		val, err := strconv.Atoi(minStr)

		if err != nil {
			c.log.Warn("redis data corrupted: invalid integer", slog.String("date", dateKey), slog.String("value", minStr))
			continue
		}

		graph[dateKey] = val
	}

	return graph, nil

}

func (c *cache) SetGraph(ctx context.Context, data map[string]int, ttl time.Duration, habitID, year string) {
	key := buildKey(habitID, year)

	pipe := c.redis.Pipeline()

	pipe.HSet(ctx, key, data)
	pipe.Expire(ctx, key, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		c.log.Error("error in repository : redis error : set graph error", slog.Any("err:", err))
	}
}

func (c *cache) DeleteGraph(ctx context.Context, id, year string) error {
	key := buildKey(id, year)

	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		c.log.Error("error in repository : redis error : delete graph error", slog.Any("err:", err))
	}

	return err
}

func (c *cache) DeleteHabit(ctx context.Context, id string) error {
	patern := buildKey(id, "")
	var errsCount int = 0

	iter := c.redis.Scan(ctx, 0, patern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		err := c.redis.Del(ctx, key).Err()
		if err != nil {
			errsCount++
			c.log.Error("error in repository : redis error : delete habbit -> delete graph error", slog.Any("err:", err))
		}
	}

	if errsCount != 0 {
		return errors.UnknownErr()
	}

	err := iter.Err()
	if err != nil {
		c.log.Error("error in repository : redis error : delete habbit -> scan error", slog.Any("err:", err))
	}

	return err
}
