package service_test

import (
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
)

const (
	LOGGER_LVL = slog.LevelError
)

type mockRepository struct {
	GetGraphFunc       func(ctx context.Context, habitID string, year int) (map[string]int, error)
	UpsertActivityFunc func(ctx context.Context, habitID string, userData domain.Track) error
	CreateHabitFunc    func(ctx context.Context, habit domain.Habit) (*domain.Habit, error)
	GetHabitsFunc      func(ctx context.Context) ([]domain.Habit, error)
	UpdateHabitFunc    func(ctx context.Context, id string, title, desc *string) (*domain.Habit, error)
	DeleteHabitFunc    func(ctx context.Context, id string) error
}

func (m *mockRepository) GetGraph(ctx context.Context, habitID string, year int) (map[string]int, error) {
	if m.GetGraphFunc != nil {
		return m.GetGraphFunc(ctx, habitID, year)
	}

	return nil, nil
}

func (m *mockRepository) UpsertActivity(ctx context.Context, habitID string, userData domain.Track) error {
	if m.UpsertActivityFunc != nil {
		return m.UpsertActivityFunc(ctx, habitID, userData)
	}

	return nil
}

func (m *mockRepository) CreateHabit(ctx context.Context, habit domain.Habit) (*domain.Habit, error) {
	if m.CreateHabitFunc != nil {
		return m.CreateHabitFunc(ctx, habit)
	}

	return nil, nil
}

func (m *mockRepository) GetHabits(ctx context.Context) ([]domain.Habit, error) {
	if m.GetHabitsFunc != nil {
		return m.GetHabitsFunc(ctx)
	}

	return nil, nil
}

func (m *mockRepository) UpdateHabit(ctx context.Context, id string, title, desc *string) (*domain.Habit, error) {
	if m.UpdateHabitFunc != nil {
		return m.UpdateHabitFunc(ctx, id, title, desc)
	}

	return nil, nil
}

func (m *mockRepository) DeleteHabit(ctx context.Context, id string) error {
	if m.DeleteHabitFunc != nil {
		return m.DeleteHabitFunc(ctx, id)
	}

	return nil
}

type mockCache struct {
	GetGraphFunc    func(ctx context.Context, habitID, year string) (map[string]int, error)
	SetGraphFunc    func(ctx context.Context, data map[string]int, ttl time.Duration, habitID, year string)
	DeleteHabitFunc func(ctx context.Context, id string) error
	DeleteGraphFunc func(ctx context.Context, id, year string) error
}

func (m *mockCache) GetGraph(ctx context.Context, habitID, year string) (map[string]int, error) {
	if m.GetGraphFunc != nil {
		return m.GetGraphFunc(ctx, habitID, year)
	}

	return nil, nil
}

func (m *mockCache) SetGraph(ctx context.Context, data map[string]int, ttl time.Duration, habitID, year string) {
	if m.SetGraphFunc != nil {
		m.SetGraphFunc(ctx, data, ttl, habitID, year)
	}

}

func (m *mockCache) DeleteHabit(ctx context.Context, id string) error {
	if m.DeleteHabitFunc != nil {
		return m.DeleteHabitFunc(ctx, id)
	}

	return nil
}

func (m *mockCache) DeleteGraph(ctx context.Context, id, year string) error {
	if m.DeleteGraphFunc != nil {
		return m.DeleteGraphFunc(ctx, id, year)
	}

	return nil
}

type testLogWriter struct {
	t *testing.T
}

func (w *testLogWriter) Write(p []byte) (n int, err error) {
	w.t.Log(strings.TrimRight(string(p), "\n"))
	return len(p), nil
}

func testLogger(t *testing.T) *slog.Logger {
	handler := slog.NewTextHandler(&testLogWriter{t: t}, &slog.HandlerOptions{
		Level: LOGGER_LVL,
	})
	return slog.New(handler)
}
