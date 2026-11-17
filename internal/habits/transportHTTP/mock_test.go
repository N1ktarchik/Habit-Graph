package transportHTTP_test

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
)

const loggerLvl = slog.LevelError

type mockService struct {
	GetGraphFunc    func(ctx context.Context, id, year string) (map[string]int, error)
	SaveDataFunc    func(ctx context.Context, id string, userData domain.Track) error
	CreateHabitFunc func(ctx context.Context, habit domain.Habit) (*domain.Habit, error)
	GetHabitsFunc   func(ctx context.Context) ([]domain.Habit, error)
	DeleteHabitFunc func(ctx context.Context, id string) error
	UpdateHabitFunc func(ctx context.Context, id string, habitInput domain.UpdateHabitInput) (*domain.Habit, error)
}

func (m *mockService) GetGraph(ctx context.Context, id, year string) (map[string]int, error) {
	if m.GetGraphFunc != nil {
		return m.GetGraphFunc(ctx, id, year)
	}
	return nil, nil
}

func (m *mockService) SaveData(ctx context.Context, id string, userData domain.Track) error {
	if m.SaveDataFunc != nil {
		return m.SaveDataFunc(ctx, id, userData)
	}
	return nil
}

func (m *mockService) CreateHabit(ctx context.Context, habit domain.Habit) (*domain.Habit, error) {
	if m.CreateHabitFunc != nil {
		return m.CreateHabitFunc(ctx, habit)
	}
	return nil, nil
}

func (m *mockService) GetHabits(ctx context.Context) ([]domain.Habit, error) {
	if m.GetHabitsFunc != nil {
		return m.GetHabitsFunc(ctx)
	}
	return nil, nil
}

func (m *mockService) DeleteHabit(ctx context.Context, id string) error {
	if m.DeleteHabitFunc != nil {
		return m.DeleteHabitFunc(ctx, id)
	}
	return nil
}

func (m *mockService) UpdateHabit(ctx context.Context, id string, habitInput domain.UpdateHabitInput) (*domain.Habit, error) {
	if m.UpdateHabitFunc != nil {
		return m.UpdateHabitFunc(ctx, id, habitInput)
	}
	return nil, nil
}

type testLogWriter struct {
	t *testing.T
}

func (w *testLogWriter) Write(p []byte) (int, error) {
	w.t.Log(strings.TrimRight(string(p), "\n"))
	return len(p), nil
}

func testLogger(t *testing.T) *slog.Logger {
	t.Helper()
	handler := slog.NewTextHandler(&testLogWriter{t: t}, &slog.HandlerOptions{
		Level: loggerLvl,
	})
	return slog.New(handler)
}
