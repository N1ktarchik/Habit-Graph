package domain

type Habit struct {
	ID          string
	Title       string
	Description string
	CreatedAt   string
}

type UpdateHabitInput struct {
	Title       *string
	Description *string
}
