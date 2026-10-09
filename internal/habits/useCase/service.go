package useCase

import (
	"context"
	"log/slog"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
)

type service struct {
	log        *slog.Logger
	repository repository
}

type repository interface {
	GetGraphFromRedis(ctx context.Context) (map[string]string, error)
	SaveDataToRedis(ctx context.Context, userData *domain.Track) error
}

func NewService(log *slog.Logger, repo repository) *service {
	return &service{
		log:        log,
		repository: repo,
	}
}
