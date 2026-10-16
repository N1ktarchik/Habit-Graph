package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
	"github.com/N1ktarchik/habbit-graph/internal/core/errors"
)

func (s *service) CreateHabit(ctx context.Context, habit domain.Habit) (*domain.Habit, error) {

	if !domain.IsValidTitle(habit.Title) {
		s.log.Debug("ERR: BAD DATA: invalid title")
		return nil, errors.BadTitle()
	}

	if !domain.IsValidDesc(habit.Description) {
		s.log.Debug("ERR: BAD DATA: invalid description")
		return nil, errors.LongDescription()
	}

	return s.repository.CreateHabit(ctx, habit)

}

func (s *service) GetHabits(ctx context.Context) ([]domain.Habit, error) {
	return s.repository.GetHabits(ctx)
}

func (s *service) UpdateHabit(ctx context.Context, id string, habitInput domain.UpdateHabitInput) (*domain.Habit, error) {

	if habitInput.Title == nil && habitInput.Description == nil {
		s.log.Debug("ERR: BAD DATA: nothing to update")
		return nil, errors.NothingToUpdate()
	}

	if habitInput.Title != nil && !domain.IsValidTitle(*habitInput.Title) {
		s.log.Debug("ERR: BAD DATA: invalid title")
		return nil, errors.BadTitle()
	}

	if habitInput.Description != nil && !domain.IsValidDesc(*habitInput.Description) {
		s.log.Debug("ERR: BAD DATA: invalid description")
		return nil, errors.LongDescription()
	}

	return s.repository.UpdateHabit(ctx, id, habitInput.Title, habitInput.Description)

}

func (s *service) DeleteHabit(ctx context.Context, id string) error {

	err := s.repository.DeleteHabit(ctx, id)
	if err != nil {
		return err
	}

	if err := s.cache.DeleteHabit(ctx, id); err != nil {
		s.log.Warn("redis failed on DeleteHabit, starting background retry", slog.Any("error", err))

		bgCtx := context.WithoutCancel(ctx)

		go func() {
			for range 5 {
				time.Sleep(2 * time.Second)
				if err := s.cache.DeleteHabit(bgCtx, id); err == nil {
					s.log.Info("habit cache deleted successfully on retry")
					return
				}
			}
			s.log.Error("failed to delete habit cache after 5 retries (memory leak risk)", slog.String("habitID", id))
		}()

	}

	return nil
}
