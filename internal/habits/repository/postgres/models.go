package postgres

import (
	"time"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
)

type modelHabit struct {
	ID          string
	Title       string
	Description string
	CreatedAt   time.Time
}

func (m *modelHabit) toDomain() domain.Habit {
	return domain.Habit{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}
