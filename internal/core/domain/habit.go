package domain

type Habit struct {
	Id          string
	Title       string
	Description string
	Created_at  string
}

type UpdateHabitInput struct {
	Title       *string
	Description *string
}
