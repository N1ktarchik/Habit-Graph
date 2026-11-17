//go:build integration

package redis_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"github.com/N1ktarchik/habit-graph/internal/habits/repository/redis"
)

func setupRedisContainer(t *testing.T) *goredis.Client {
	t.Helper()
	ctx := context.Background()

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(redisContainer)
	})

	connStr, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	opts, err := goredis.ParseURL(connStr)
	require.NoError(t, err)

	client := goredis.NewClient(opts)
	require.NoError(t, client.Ping(ctx).Err())

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}

func TestRedisCache_Integration(t *testing.T) {
	client := setupRedisContainer(t)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	cache := redis.NewCache(log, client)
	ctx := context.Background()

	const (
		habitID      = "550e8400-e29b-41d4-a716-446655440000"
		otherHabitID = "660e8400-e29b-41d4-a716-446655440000"
	)

	t.Run("GetGraph cache miss returns nil without error", func(t *testing.T) {
		graph, err := cache.GetGraph(ctx, habitID, "2026")
		require.NoError(t, err)
		assert.Nil(t, graph)
	})

	t.Run("SetGraph stores hash with TTL and GetGraph retrieves it", func(t *testing.T) {
		expected := map[string]int{
			"2026-11-15": 45,
			"2026-11-16": 90,
		}

		cache.SetGraph(ctx, expected, 10*time.Minute, habitID, "2026")

		got, err := cache.GetGraph(ctx, habitID, "2026")
		require.NoError(t, err)
		assert.Equal(t, expected, got)

		ttl, err := client.TTL(ctx, "habit:"+habitID+":graph:2026").Result()
		require.NoError(t, err)
		assert.Greater(t, ttl, 9*time.Minute)
	})

	t.Run("GetGraph skips corrupted non-integer values gracefully", func(t *testing.T) {
		key := "habit:" + habitID + ":graph:2028"
		err := client.HSet(ctx, key, map[string]any{
			"2028-01-01": "60",
			"2028-01-02": "corrupted_value",
		}).Err()
		require.NoError(t, err)

		got, err := cache.GetGraph(ctx, habitID, "2028")
		require.NoError(t, err)
		assert.Equal(t, map[string]int{"2028-01-01": 60}, got)
	})

	t.Run("DeleteGraph removes only specified year", func(t *testing.T) {
		cache.SetGraph(ctx, map[string]int{"2026-01-01": 15}, time.Hour, habitID, "2026")
		cache.SetGraph(ctx, map[string]int{"2027-01-01": 30}, time.Hour, habitID, "2027")

		err := cache.DeleteGraph(ctx, habitID, "2026")
		require.NoError(t, err)

		got2026, err := cache.GetGraph(ctx, habitID, "2026")
		require.NoError(t, err)
		assert.Nil(t, got2026)

		got2027, err := cache.GetGraph(ctx, habitID, "2027")
		require.NoError(t, err)
		assert.Equal(t, map[string]int{"2027-01-01": 30}, got2027)
	})

	t.Run("DeleteHabit scans and deletes all years for target habit only", func(t *testing.T) {
		cache.SetGraph(ctx, map[string]int{"2026-02-01": 10}, time.Hour, habitID, "2026")
		cache.SetGraph(ctx, map[string]int{"2027-02-01": 20}, time.Hour, habitID, "2027")
		cache.SetGraph(ctx, map[string]int{"2026-02-01": 99}, time.Hour, otherHabitID, "2026")

		err := cache.DeleteHabit(ctx, habitID)
		require.NoError(t, err)

		g1, err := cache.GetGraph(ctx, habitID, "2026")
		require.NoError(t, err)
		assert.Nil(t, g1)

		g2, err := cache.GetGraph(ctx, habitID, "2027")
		require.NoError(t, err)
		assert.Nil(t, g2)

		otherGraph, err := cache.GetGraph(ctx, otherHabitID, "2026")
		require.NoError(t, err)
		assert.Equal(t, map[string]int{"2026-02-01": 99}, otherGraph)
	})
}
