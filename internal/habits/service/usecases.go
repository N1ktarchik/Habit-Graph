package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
)

func (s *service) GetGraph(ctx context.Context, id, year string) (map[string]string, error) {
	if !domain.IsValidYear(year) {
		s.log.Debug("ERR: BAD DATA: invalid year")
		return nil, errors.BadYear()
	}

	m, err := s.cache.GetGraph(ctx, id, year)
	if err == nil && len(m) != 0 {
		return m, nil
	}

	m, err = s.repository.GetGraph(ctx, id, year)
	if err != nil {
		return nil, err
	}

	if len(m) != 0 {
		go func() {
			bgCtx := context.WithoutCancel(ctx)
			s.cache.SetGraph(bgCtx, m, 24*time.Hour, id, year)
		}()
	}

	return m, nil
}

func (s *service) SaveData(ctx context.Context, id string, userData domain.Track) error {
	if userData.Value < 0 {
		s.log.Debug("ERR: BAD DATA: invalid value")
		return errors.BadValue()
	}

	if !domain.IsValidDate(userData.Date) {
		s.log.Debug("ERR: BAD DATA: invalid date")
		return errors.BadDate()
	}

	if err := s.repository.UpsertActivity(ctx, id, userData); err != nil {
		return err
	}

	year := userData.Date[:4]
	if err := s.cache.DeleteGraph(ctx, id, year); err != nil {

		s.log.Warn("redis failed, starting background retry", slog.Any("error", err))

		go func() {
			bgCtx := context.WithoutCancel(ctx)

			for range 5 {
				time.Sleep(2 * time.Second)
				if err := s.cache.DeleteGraph(bgCtx, id, year); err == nil {
					s.log.Info("cache deleted successfully on retry")
					return
				}
			}

			s.log.Error("failed to delete cache after 5 retries", slog.String("habitID", id))
		}()
	}

	return nil

}
