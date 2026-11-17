//go:build integration

package postgres_test

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
	coreErrors "github.com/N1ktarchik/habit-graph/internal/core/errors"
	"github.com/N1ktarchik/habit-graph/internal/habits/repository/postgres"
)

func setupPostgresContainer(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("habit_graph_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("test123"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(pgContainer)
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	migrationPath := filepath.Join("..", "..", "..", "..", "migrations", "000001_init.up.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, string(sqlBytes))
	require.NoError(t, err)

	return pool
}

func TestPostgresRepository_Integration(t *testing.T) {
	pool := setupPostgresContainer(t)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := postgres.NewRepository(log, pool)
	ctx := context.Background()

	const nonExistentUUID = "00000000-0000-0000-0000-000000000000"

	t.Run("CreateHabit and GetHabits", func(t *testing.T) {
		created, err := repo.CreateHabit(ctx, domain.Habit{
			Title:       "Go Backend",
			Description: "Clean Architecture & Testing",
		})
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, "Go Backend", created.Title)
		assert.Equal(t, "Clean Architecture & Testing", created.Description)
		assert.False(t, created.CreatedAt.IsZero())

		habits, err := repo.GetHabits(ctx)
		require.NoError(t, err)
		require.Len(t, habits, 1)
		assert.Equal(t, created.ID, habits[0].ID)
	})

	t.Run("UpdateHabit with COALESCE and NotFound", func(t *testing.T) {
		created, err := repo.CreateHabit(ctx, domain.Habit{
			Title:       "Initial Title",
			Description: "Initial Desc",
		})
		require.NoError(t, err)

		newTitle := "Updated Title"
		updated, err := repo.UpdateHabit(ctx, created.ID, &newTitle, nil)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "Initial Desc", updated.Description)

		newDesc := "Updated Desc"
		updated2, err := repo.UpdateHabit(ctx, created.ID, nil, &newDesc)
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated2.Title)
		assert.Equal(t, "Updated Desc", updated2.Description)

		res, err := repo.UpdateHabit(ctx, nonExistentUUID, &newTitle, nil)
		assert.Nil(t, res)
		assert.Equal(t, coreErrors.NotFound(), err)
	})

	t.Run("UpsertActivity, GetGraph and DeleteHabit cascade", func(t *testing.T) {
		habit, err := repo.CreateHabit(ctx, domain.Habit{
			Title:       "Workout",
			Description: "Gym sessions",
		})
		require.NoError(t, err)

		err = repo.UpsertActivity(ctx, habit.ID, domain.Track{Date: "2026-05-10", Value: 30})
		require.NoError(t, err)

		err = repo.UpsertActivity(ctx, habit.ID, domain.Track{Date: "2026-05-10", Value: 45})
		require.NoError(t, err)

		err = repo.UpsertActivity(ctx, habit.ID, domain.Track{Date: "2026-05-11", Value: 20})
		require.NoError(t, err)
		err = repo.UpsertActivity(ctx, habit.ID, domain.Track{Date: "2027-01-01", Value: 60})
		require.NoError(t, err)

		graph2026, err := repo.GetGraph(ctx, habit.ID, 2026)
		require.NoError(t, err)
		assert.Equal(t, map[string]int{
			"2026-05-10": 75,
			"2026-05-11": 20,
		}, graph2026)

		err = repo.UpsertActivity(ctx, habit.ID, domain.Track{Date: "2026-05-12", Value: -10})
		assert.Error(t, err)

		err = repo.DeleteHabit(ctx, habit.ID)
		require.NoError(t, err)

		err = repo.DeleteHabit(ctx, habit.ID)
		assert.Equal(t, coreErrors.NotFound(), err)

		graphAfterDelete, err := repo.GetGraph(ctx, habit.ID, 2026)
		require.NoError(t, err)
		assert.Empty(t, graphAfterDelete)
	})
}
