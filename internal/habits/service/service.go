package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
)

type service struct {
	log        *slog.Logger
	repository repository
	cache      cache
}

type repository interface {
	GetGraph(ctx context.Context, habitID, year string) (map[string]string, error)
	UpsertActivity(ctx context.Context, habitID string, userData domain.Track) error

	CreateHabit(ctx context.Context, habit domain.Habit) (*domain.Habit, error)
	GetHabits(ctx context.Context) ([]domain.Habit, error)
	UpdateHabit(ctx context.Context, id string, title, desc *string) (*domain.Habit, error)
	DeleteHabit(ctx context.Context, id string) error
}

type cache interface {
	GetGraph(ctx context.Context, habitID, year string) (map[string]string, error)
	SetGraph(ctx context.Context, data map[string]string, ttl time.Duration, habitID, year string)

	DeleteHabit(ctx context.Context, id string) error
	DeleteGraph(ctx context.Context, id, year string) error
}

func NewService(log *slog.Logger, repo repository, cache cache) *service {
	return &service{
		log:        log,
		repository: repo,
		cache:      cache,
	}
}
