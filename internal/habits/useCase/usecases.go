package useCase

import (
	"context"
	"time"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
)

func (s *service) GetGraph(ctx context.Context) (map[string]string, error) {
	return s.repository.GetGraphFromRedis(ctx)
}

func (s *service) SaveData(ctx context.Context, userData *domain.UserData) error {
	if userData.Value < 0 {
		return errors.BadValue()
	}

	parsedDate, err := time.Parse(time.DateOnly, userData.Date)
	if err != nil {
		return errors.BadDate()
	}

	currentYear := time.Now().Year()
	if parsedDate.Year() < currentYear-1 || parsedDate.Year() > currentYear+1 {
		return errors.BadDate()
	}

	return s.repository.SaveDataToRedis(ctx, userData)

}
