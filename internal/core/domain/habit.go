package domain

import "time"

type Habit struct {
	ID          string
	Title       string
	Description string
	CreatedAt   time.Time
}

type UpdateHabitInput struct {
	Title       *string
	Description *string
}
