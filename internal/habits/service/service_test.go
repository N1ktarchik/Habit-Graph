package service_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	core_errors "github.com/N1ktarchik/habit-graph/internal/core/errors"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
	"github.com/N1ktarchik/habit-graph/internal/habits/service"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}

	return *p
}

func TestCreateHabit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input       domain.Habit
		expectedErr error
		wantDBCall  bool
	}{
		{
			name:        "valid habit",
			input:       domain.Habit{Title: "Test Habit 1", Description: "Description 1"},
			expectedErr: nil,
			wantDBCall:  true,
		},

		{
			name:        "missing title",
			input:       domain.Habit{Description: "Description 2"},
			expectedErr: core_errors.BadTitle(),
			wantDBCall:  false,
		},
		{
			name:        "space title",
			input:       domain.Habit{Title: " ", Description: "Description 3"},
			expectedErr: core_errors.BadTitle(),
			wantDBCall:  false,
		},
		{
			name:        "empty description",
			input:       domain.Habit{Title: "Test Habit 4", Description: ""},
			expectedErr: nil,
			wantDBCall:  true,
		},
		{
			name:        "too long description",
			input:       domain.Habit{Title: "Test Habit 5", Description: "This is a very long description that exceeds the maximum allowed length for a habit description. It should trigger a validation error when trying to create the habit. ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
			expectedErr: core_errors.BadDescription(),
			wantDBCall:  false,
		},
		{
			name:        "missing description",
			input:       domain.Habit{Title: "Test Habit 6"},
			expectedErr: nil,
			wantDBCall:  true,
		},
		{
			name:        "too long title",
			input:       domain.Habit{Title: "This is a very long habit title that exceeds the maximum allowed length for a habit title. It should trigger a validation error when trying to create the habit. ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", Description: "Description 7"},
			expectedErr: core_errors.BadTitle(),
			wantDBCall:  false,
		},
		{
			name:        "valid habit with client-side ID to be overridden",
			input:       domain.Habit{ID: "test-id-8", Title: "Test Habit 8", Description: "Description 8"},
			expectedErr: nil,
			wantDBCall:  true,
		},
		{
			name:        "valid habit with client-side CreatedAt",
			input:       domain.Habit{Title: "Test Habit 9", Description: "Description 9", CreatedAt: time.Now().UTC()},
			expectedErr: nil,
			wantDBCall:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log := testLogger(t)
			dbCalled := false

			r := &mockRepository{
				CreateHabitFunc: func(ctx context.Context, h domain.Habit) (*domain.Habit, error) {
					dbCalled = true
					return &domain.Habit{
						ID:          "server-generated-id",
						Title:       h.Title,
						Description: h.Description,
						CreatedAt:   time.Now(),
					}, nil
				},
			}

			srv := service.NewService(log, r, nil)

			result, err := srv.CreateHabit(context.Background(), tc.input)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Nil(t, result, "при ошибке результат должен быть nil")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.ID)
				assert.False(t, result.CreatedAt.IsZero())
				assert.Equal(t, tc.input.Title, result.Title)
				assert.Equal(t, tc.input.Description, result.Description)
			}

			assert.Equal(t, tc.wantDBCall, dbCalled, "состояние вызова репозитория не совпадает")
		})
	}
}

func TestUpdateHabit(t *testing.T) {
	t.Parallel()

	const (
		titleFromDB = "Test Update Title "
		descFromDB  = "Test Update Description "
	)

	testCases := []struct {
		name        string
		inputID     string
		inputData   domain.UpdateHabitInput
		expectedErr error
		wantDBCall  bool
	}{
		{
			name:        "valid habit",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: ptr("Test Update Habit "), Description: ptr("Description ")},
			expectedErr: nil,
			wantDBCall:  true,
		},

		{
			name:        "valid habit/only title",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: ptr("Test Update Habit "), Description: nil},
			expectedErr: nil,
			wantDBCall:  true,
		},

		{
			name:        "valid habit/only description",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: nil, Description: ptr("Description ")},
			expectedErr: nil,
			wantDBCall:  true,
		},

		{
			name:        "nil ptrs",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: nil, Description: nil},
			expectedErr: core_errors.NothingToUpdate(),
			wantDBCall:  false,
		},

		{
			name:        "space title and nil desc",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: ptr(" "), Description: nil},
			expectedErr: core_errors.BadTitle(),
			wantDBCall:  false,
		},

		{
			name:        "space title and space desc",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: ptr(" "), Description: ptr(" ")},
			expectedErr: core_errors.BadTitle(),
			wantDBCall:  false,
		},

		{
			name:        "too long title",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: ptr("This is a very long habit title that exceeds the maximum allowed length for a habit title. It should trigger a validation error when trying to create the habit. ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), Description: nil},
			expectedErr: core_errors.BadTitle(),
			wantDBCall:  false,
		},

		{
			name:        "too long description",
			inputID:     "habit-ID",
			inputData:   domain.UpdateHabitInput{Title: nil, Description: ptr("This is a very long description that exceeds the maximum allowed length for a habit description. It should trigger a validation error when trying to create the habit. ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			expectedErr: core_errors.BadDescription(),
			wantDBCall:  false,
		},

		{
			name:        "id not found",
			inputID:     " ",
			inputData:   domain.UpdateHabitInput{Title: ptr("Test Update Habit "), Description: ptr("Description ")},
			expectedErr: core_errors.UnknownErr(),
			wantDBCall:  true,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log := testLogger(t)
			dbCalled := false

			r := &mockRepository{
				UpdateHabitFunc: func(ctx context.Context, id string, title *string, desc *string) (*domain.Habit, error) {

					var tt, dsc = titleFromDB, descFromDB

					dbCalled = true

					if id == " " {
						return nil, core_errors.UnknownErr()
					}

					if title != nil {
						tt = *title
					}

					if desc != nil {
						dsc = *desc
					}

					return &domain.Habit{
						ID:          id,
						Title:       tt,
						Description: dsc,
						CreatedAt:   time.Now(),
					}, nil
				},
			}

			srv := service.NewService(log, r, nil)

			result, err := srv.UpdateHabit(context.Background(), tc.inputID, tc.inputData)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Nil(t, result, "при ошибке результат должен быть nil")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.inputID, result.ID)

				if tc.inputData.Title != nil {
					assert.Equal(t, deref(tc.inputData.Title), result.Title)
				} else {
					assert.Equal(t, titleFromDB, result.Title)
				}

				if tc.inputData.Description != nil {
					assert.Equal(t, deref(tc.inputData.Description), result.Description)
				} else {
					assert.Equal(t, descFromDB, result.Description)
				}

				assert.False(t, result.CreatedAt.IsZero())
			}

			assert.Equal(t, tc.wantDBCall, dbCalled, "состояние вызова репозитория не совпадает")

		})
	}
}

func TestDeleteHabit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		inputId          string
		expectedErr      error
		wantPostgresCall bool
		wantRedisCall    bool
		wantRetries      bool
		expectedRetries  int32
	}{
		{
			name:             "Valid id",
			inputId:          "habit-id",
			expectedErr:      nil,
			wantPostgresCall: true,
			wantRedisCall:    true,
			wantRetries:      false,
		},
		{
			name:             "postgres err/habit not found",
			inputId:          " ",
			expectedErr:      core_errors.NotFound(),
			wantPostgresCall: true,
			wantRedisCall:    false,
			wantRetries:      false,
		},
		{
			name:             "redis err/triggers retries",
			inputId:          "habit-id-retry",
			expectedErr:      nil,
			wantPostgresCall: true,
			wantRedisCall:    true,
			wantRetries:      true,
			expectedRetries:  6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log := testLogger(t)
			postgresCalled := false

			var cacheCalls atomic.Int32

			r := &mockRepository{
				DeleteHabitFunc: func(ctx context.Context, id string) error {
					postgresCalled = true
					if id == " " {
						return core_errors.NotFound()
					}
					return nil
				},
			}

			c := &mockCache{
				DeleteHabitFunc: func(ctx context.Context, id string) error {
					cacheCalls.Add(1)

					if tc.wantRetries {
						return core_errors.UnknownErr()
					}
					return nil
				},
			}

			srv := service.NewService(log, r, c)

			err := srv.DeleteHabit(context.Background(), tc.inputId)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantPostgresCall, postgresCalled, "Postgres call mismatch")

			if tc.wantRetries {
				assert.Eventually(t, func() bool {
					return cacheCalls.Load() == tc.expectedRetries
				}, 12*time.Second, 1*time.Second, "Фоновая горутина не выполнила нужного количества вызовов")
			} else {
				if tc.wantRedisCall {
					assert.Equal(t, int32(1), cacheCalls.Load(), "Ожидался 1 вызов Redis")
				} else {
					assert.Equal(t, int32(0), cacheCalls.Load(), "Ожидалось 0 вызовов Redis")
				}
			}
		})
	}
}
