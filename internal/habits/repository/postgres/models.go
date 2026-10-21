package postgres

import (
	"time"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
)

type modelHabit struct {
	ID          string
	Title       string
	Description string
	CreatedAt   time.Time
}

type modelTrack struct {
	Date  string
	Value int
}

func (m *modelHabit) toDomain() domain.Habit {
	return domain.Habit{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}

func (m *modelTrack) toDomain() domain.Track {
	return domain.Track{
		Date:  m.Date,
		Value: m.Value,
	}
}
