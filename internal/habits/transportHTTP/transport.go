package transportHTTP

import (
	"context"
	"log/slog"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
)

type transportHTTP struct {
	log  *slog.Logger
	serv service
}

type service interface {
	GetGraph(ctx context.Context) (map[string]string, error)
	SaveData(ctx context.Context, userData *domain.UserData) error
}

func NewTransportHTTP(log *slog.Logger, serv service) *transportHTTP {
	return &transportHTTP{
		log:  log,
		serv: serv,
	}
}
