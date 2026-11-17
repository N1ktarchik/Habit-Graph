package service_test

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	core_errors "github.com/N1ktarchik/habit-graph/internal/core/errors"

	"github.com/stretchr/testify/assert"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
	"github.com/N1ktarchik/habit-graph/internal/habits/service"
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

func TestGetGraph(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		inputId          string
		inputYear        string
		wantPostgresCall bool
		wantPostgresErr  bool
		wantRedisGetCall bool
		wantRedisSetCall bool
		expectedErr      error
	}{
		{
			name:             "Valid input data (minimum year) and Redis cache hit",
			inputId:          "habit-id",
			inputYear:        "2026",
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisGetCall: true,
			wantRedisSetCall: false,
			expectedErr:      nil,
		},
		{
			name:             "Valid input data (current year) and Redis cache hit",
			inputId:          "habit-id",
			inputYear:        strconv.Itoa(time.Now().Year()),
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisGetCall: true,
			wantRedisSetCall: false,
			expectedErr:      nil,
		},
		{
			name:             "Valid input data and Postgres cache hit",
			inputId:          "habit-id",
			inputYear:        "2026",
			wantPostgresCall: true,
			wantPostgresErr:  false,
			wantRedisGetCall: true,
			wantRedisSetCall: true,
			expectedErr:      nil,
		},
		{
			name:             "Invalid input year : old year",
			inputId:          "habit-id",
			inputYear:        "2020",
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisGetCall: false,
			wantRedisSetCall: false,
			expectedErr:      core_errors.BadYear(),
		},

		{
			name:             "Invalid input year : future year",
			inputId:          "habit-id",
			inputYear:        "1010020",
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisGetCall: false,
			wantRedisSetCall: false,
			expectedErr:      core_errors.BadYear(),
		},

		{
			name:             "Postgres returns error",
			inputId:          "habit-id",
			inputYear:        "2026",
			wantPostgresCall: true,
			wantPostgresErr:  true,
			wantRedisGetCall: true,
			wantRedisSetCall: false,
			expectedErr:      core_errors.UnknownErr(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			postgresCalled := false
			redisGetCalled := false
			var redisSetCalled atomic.Bool

			log := testLogger(t)

			r := &mockRepository{
				GetGraphFunc: func(ctx context.Context, habitID string, year int) (map[string]int, error) {
					postgresCalled = true
					if tc.wantPostgresErr {
						return nil, core_errors.UnknownErr()
					}
					return map[string]int{"2023-01-01": 1}, nil
				},
			}

			c := &mockCache{
				GetGraphFunc: func(ctx context.Context, habitID, year string) (map[string]int, error) {
					redisGetCalled = true
					if tc.wantPostgresCall {
						return nil, nil
					}

					return map[string]int{"2023-01-01": 1}, nil
				},

				SetGraphFunc: func(ctx context.Context, data map[string]int, ttl time.Duration, habitID, year string) {
					redisSetCalled.Store(true)
				},
			}

			srv := service.NewService(log, r, c)

			result, err := srv.GetGraph(context.Background(), tc.inputId, tc.inputYear)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Nil(t, result, "при ошибке результат должен быть nil")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, len(result), 1)
			}

			assert.Equal(t, tc.wantPostgresCall, postgresCalled, "состояние вызова Postgres не совпадает")
			assert.Equal(t, tc.wantRedisGetCall, redisGetCalled, "состояние вызова Redis Get не совпадает")

			if tc.wantRedisSetCall {
				assert.Eventually(
					t,
					func() bool { return redisSetCalled.Load() },
					2*time.Second,
					10*time.Millisecond,
					"состояние вызова Redis Set не совпадает (горутина не отработала)")
			} else {
				assert.False(t, redisSetCalled.Load(), "запись в Redis не ожидалась")
			}

		})
	}
}

func TestSaveData(t *testing.T) {
	t.Parallel()

	today := time.Now().Format(time.DateOnly)
	futureDate := time.Now().Add(72 * time.Hour).Format(time.DateOnly)

	testCases := []struct {
		name               string
		inputID            string
		inputData          domain.Track
		wantPostgresCall   bool
		wantPostgresErr    bool
		wantRedisCalls     int32
		redisFailUntilCall int32
		expectedErr        error
	}{
		{
			name:             "valid track data",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: today, Value: 5},
			wantPostgresCall: true,
			wantPostgresErr:  false,
			wantRedisCalls:   1,
			expectedErr:      nil,
		},
		{
			name:             "valid track data with zero value",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: "2026-01-01", Value: 0},
			wantPostgresCall: true,
			wantPostgresErr:  false,
			wantRedisCalls:   1,
			expectedErr:      nil,
		},
		{
			name:             "negative value",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: today, Value: -1},
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisCalls:   0,
			expectedErr:      core_errors.BadValue(),
		},
		{
			name:             "invalid date format",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: "16-11-2026", Value: 1},
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisCalls:   0,
			expectedErr:      core_errors.BadDate(),
		},
		{
			name:             "date year before 2026",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: "2025-12-31", Value: 1},
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisCalls:   0,
			expectedErr:      core_errors.BadDate(),
		},
		{
			name:             "future date beyond 24h limit",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: futureDate, Value: 1},
			wantPostgresCall: false,
			wantPostgresErr:  false,
			wantRedisCalls:   0,
			expectedErr:      core_errors.BadDate(),
		},
		{
			name:             "postgres error",
			inputID:          "habit-id",
			inputData:        domain.Track{Date: today, Value: 1},
			wantPostgresCall: true,
			wantPostgresErr:  true,
			wantRedisCalls:   0,
			expectedErr:      core_errors.UnknownErr(),
		},
		{
			name:               "redis fails once and succeeds on 1st retry",
			inputID:            "habit-id",
			inputData:          domain.Track{Date: today, Value: 1},
			wantPostgresCall:   true,
			wantPostgresErr:    false,
			wantRedisCalls:     2,
			redisFailUntilCall: 1,
			expectedErr:        nil,
		},
		{
			name:               "redis fails all 5 retries",
			inputID:            "habit-id",
			inputData:          domain.Track{Date: today, Value: 1},
			wantPostgresCall:   true,
			wantPostgresErr:    false,
			wantRedisCalls:     6,
			redisFailUntilCall: 6,
			expectedErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log := testLogger(t)
			postgresCalled := false
			var redisCalls atomic.Int32

			r := &mockRepository{
				UpsertActivityFunc: func(ctx context.Context, id string, data domain.Track) error {
					postgresCalled = true
					assert.Equal(t, tc.inputID, id)
					assert.Equal(t, tc.inputData, data)

					if tc.wantPostgresErr {
						return core_errors.UnknownErr()
					}
					return nil
				},
			}

			c := &mockCache{
				DeleteGraphFunc: func(ctx context.Context, habitID, year string) error {
					currentCall := redisCalls.Add(1)
					assert.Equal(t, tc.inputID, habitID)
					assert.Equal(t, tc.inputData.Date[:4], year)

					if currentCall <= tc.redisFailUntilCall {
						return errors.New("redis connection error")
					}
					return nil
				},
			}

			srv := service.NewService(log, r, c)

			err := srv.SaveData(context.Background(), tc.inputID, tc.inputData)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantPostgresCall, postgresCalled, "состояние вызова Postgres не совпадает")

			if tc.wantRedisCalls > 1 {
				assert.Eventually(t, func() bool {
					return redisCalls.Load() == tc.wantRedisCalls
				}, 12*time.Second, 100*time.Millisecond, "количество вызовов Redis с учетом ретраев не совпадает")
			} else {
				assert.Equal(t, tc.wantRedisCalls, redisCalls.Load(), "количество вызовов Redis не совпадает")
			}
		})
	}
}
