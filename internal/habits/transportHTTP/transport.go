package transportHTTP

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
	"github.com/N1ktarchik/habbit-graph/internal/core/response"
)

type transportHTTP struct {
	log     *slog.Logger
	service service
}

type service interface {
	GetGraph(ctx context.Context, id, year string) (map[string]string, error)
	SaveData(ctx context.Context, id string, userData domain.Track) error

	CreateHabit(ctx context.Context, habit domain.Habit) (*domain.Habit, error)
	GetHabits(ctx context.Context) ([]domain.Habit, error)
	DeleteHabit(ctx context.Context, id string) error
	UpdateHabit(ctx context.Context, id string, habitInput *domain.UpdateHabitInput) (*domain.Habit, error)
}

func NewTransportHTTP(log *slog.Logger, serv service) *transportHTTP {
	return &transportHTTP{
		log:     log,
		service: serv,
	}
}

func (t *transportHTTP) sendError(w http.ResponseWriter, err error) {
	if respErr := response.ResponseWithError(w, err); respErr != nil {
		t.log.Warn("error in response function", slog.Any("err", respErr))
	}
}
