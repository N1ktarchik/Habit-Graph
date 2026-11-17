package transportHTTP_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/N1ktarchik/habit-graph/internal/core/domain"
	coreErrors "github.com/N1ktarchik/habit-graph/internal/core/errors"
	transport "github.com/N1ktarchik/habit-graph/internal/habits/transportHTTP"
)

const validUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestCreateHabit(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		reqBody        string
		mockReturn     *domain.Habit
		mockErr        error
		wantSrvCall    bool
		expectedStatus int
		expectedInBody string
	}{
		{
			name:    "valid habit -> 201 Created",
			reqBody: `{"title": "Go Backend", "description": "Write unit tests"}`,
			mockReturn: &domain.Habit{
				ID:          validUUID,
				Title:       "Go Backend",
				Description: "Write unit tests",
				CreatedAt:   fixedTime,
			},
			mockErr:        nil,
			wantSrvCall:    true,
			expectedStatus: http.StatusCreated,
			expectedInBody: `"id":"` + validUUID + `"`,
		},
		{
			name:           "invalid JSON -> 400 Bad Request",
			reqBody:        `{"title": "broken json"`,
			mockReturn:     nil,
			mockErr:        nil,
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid JSON"`,
		},
		{
			name:           "service validation error -> 400 Bad Request",
			reqBody:        `{"title": "a"}`,
			mockReturn:     nil,
			mockErr:        coreErrors.BadTitle(),
			wantSrvCall:    true,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid title"`,
		},
		{
			name:           "unexpected repository error -> 500 Internal Server Error",
			reqBody:        `{"title": "Go Backend"}`,
			mockReturn:     nil,
			mockErr:        errors.New("db connection lost"),
			wantSrvCall:    true,
			expectedStatus: http.StatusInternalServerError,
			expectedInBody: `"internal service error"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srvCalled := false
			srv := &mockService{
				CreateHabitFunc: func(ctx context.Context, habit domain.Habit) (*domain.Habit, error) {
					srvCalled = true
					return tc.mockReturn, tc.mockErr
				},
			}

			h := transport.NewTransportHTTP(testLogger(t), srv)

			req := httptest.NewRequest(http.MethodPost, "/api/habits", bytes.NewBufferString(tc.reqBody))
			rec := httptest.NewRecorder()

			h.CreateHabit(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.wantSrvCall, srvCalled)
			if tc.expectedInBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedInBody)
			}
		})
	}
}

func TestGetHabits(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		mockReturn     []domain.Habit
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "returns list of habits -> 200 OK",
			mockReturn: []domain.Habit{
				{ID: validUUID, Title: "Habit 1", Description: "Desc 1", CreatedAt: fixedTime},
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Habit 1"`,
		},
		{
			name:           "empty list returns empty JSON array [] -> 200 OK",
			mockReturn:     []domain.Habit{},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "[]",
		},
		{
			name:           "service error -> 500 Internal Server Error",
			mockReturn:     nil,
			mockErr:        coreErrors.UnknownErr(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"internal service error"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := &mockService{
				GetHabitsFunc: func(ctx context.Context) ([]domain.Habit, error) {
					return tc.mockReturn, tc.mockErr
				},
			}

			h := transport.NewTransportHTTP(testLogger(t), srv)

			req := httptest.NewRequest(http.MethodGet, "/api/habits", nil)
			rec := httptest.NewRecorder()

			h.GetHabits(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		})
	}
}

func TestUpdateHabit(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		pathID         string
		reqBody        string
		mockReturn     *domain.Habit
		mockErr        error
		wantSrvCall    bool
		expectedStatus int
		expectedInBody string
	}{
		{
			name:    "valid update -> 200 OK",
			pathID:  validUUID,
			reqBody: `{"title": "Updated Title"}`,
			mockReturn: &domain.Habit{
				ID:          validUUID,
				Title:       "Updated Title",
				Description: "Old Desc",
				CreatedAt:   fixedTime,
			},
			mockErr:        nil,
			wantSrvCall:    true,
			expectedStatus: http.StatusOK,
			expectedInBody: `"title":"Updated Title"`,
		},
		{
			name:           "invalid UUID in path -> 400 Bad Request",
			pathID:         "not-a-uuid",
			reqBody:        `{"title": "Updated Title"}`,
			mockReturn:     nil,
			mockErr:        nil,
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid ID"`,
		},
		{
			name:           "invalid JSON body -> 400 Bad Request",
			pathID:         validUUID,
			reqBody:        `{broken}`,
			mockReturn:     nil,
			mockErr:        nil,
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid JSON"`,
		},
		{
			name:           "habit not found -> 404 Not Found",
			pathID:         validUUID,
			reqBody:        `{"title": "Updated Title"}`,
			mockReturn:     nil,
			mockErr:        coreErrors.NotFound(),
			wantSrvCall:    true,
			expectedStatus: http.StatusNotFound,
			expectedInBody: `"data not found"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srvCalled := false
			srv := &mockService{
				UpdateHabitFunc: func(ctx context.Context, id string, input domain.UpdateHabitInput) (*domain.Habit, error) {
					srvCalled = true
					assert.Equal(t, tc.pathID, id)
					return tc.mockReturn, tc.mockErr
				},
			}

			h := transport.NewTransportHTTP(testLogger(t), srv)

			req := httptest.NewRequest(http.MethodPatch, "/api/habits/"+tc.pathID, bytes.NewBufferString(tc.reqBody))
			req.SetPathValue("id", tc.pathID)
			rec := httptest.NewRecorder()

			h.UpdateHabit(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.wantSrvCall, srvCalled)
			assert.Contains(t, rec.Body.String(), tc.expectedInBody)
		})
	}
}

func TestDeleteHabit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		pathID         string
		mockErr        error
		wantSrvCall    bool
		expectedStatus int
	}{
		{
			name:           "valid delete -> 204 No Content",
			pathID:         validUUID,
			mockErr:        nil,
			wantSrvCall:    true,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid UUID -> 400 Bad Request",
			pathID:         "invalid-uuid",
			mockErr:        nil,
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "habit not found -> 404 Not Found",
			pathID:         validUUID,
			mockErr:        coreErrors.NotFound(),
			wantSrvCall:    true,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srvCalled := false
			srv := &mockService{
				DeleteHabitFunc: func(ctx context.Context, id string) error {
					srvCalled = true
					assert.Equal(t, tc.pathID, id)
					return tc.mockErr
				},
			}

			h := transport.NewTransportHTTP(testLogger(t), srv)

			req := httptest.NewRequest(http.MethodDelete, "/api/habits/"+tc.pathID, nil)
			req.SetPathValue("id", tc.pathID)
			rec := httptest.NewRecorder()

			h.DeleteHabit(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.wantSrvCall, srvCalled)
		})
	}
}

func TestGetGraph(t *testing.T) {
	t.Parallel()

	currentYear := strconv.Itoa(time.Now().Year())

	testCases := []struct {
		name           string
		pathID         string
		queryYear      string
		expectedYear   string
		mockReturn     map[string]int
		mockErr        error
		wantSrvCall    bool
		expectedStatus int
		expectedInBody string
	}{
		{
			name:           "valid request with explicit year -> 200 OK",
			pathID:         validUUID,
			queryYear:      "2026",
			expectedYear:   "2026",
			mockReturn:     map[string]int{"2026-11-16": 45},
			mockErr:        nil,
			wantSrvCall:    true,
			expectedStatus: http.StatusOK,
			expectedInBody: `"2026-11-16":45`,
		},
		{
			name:           "missing year query param defaults to current year -> 200 OK",
			pathID:         validUUID,
			queryYear:      "",
			expectedYear:   currentYear,
			mockReturn:     map[string]int{"2026-01-01": 10},
			mockErr:        nil,
			wantSrvCall:    true,
			expectedStatus: http.StatusOK,
			expectedInBody: `"2026-01-01":10`,
		},
		{
			name:           "invalid UUID -> 400 Bad Request",
			pathID:         "bad-id",
			queryYear:      "2026",
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid ID"`,
		},
		{
			name:           "invalid year from service -> 400 Bad Request",
			pathID:         validUUID,
			queryYear:      "2020",
			expectedYear:   "2020",
			mockReturn:     nil,
			mockErr:        coreErrors.BadYear(),
			wantSrvCall:    true,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid year"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srvCalled := false
			srv := &mockService{
				GetGraphFunc: func(ctx context.Context, id, year string) (map[string]int, error) {
					srvCalled = true
					assert.Equal(t, tc.pathID, id)
					assert.Equal(t, tc.expectedYear, year)
					return tc.mockReturn, tc.mockErr
				},
			}

			h := transport.NewTransportHTTP(testLogger(t), srv)

			url := "/api/habits/" + tc.pathID + "/graph"
			if tc.queryYear != "" {
				url += "?year=" + tc.queryYear
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.SetPathValue("id", tc.pathID)
			rec := httptest.NewRecorder()

			h.GetGraph(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.wantSrvCall, srvCalled)
			assert.Contains(t, rec.Body.String(), tc.expectedInBody)
		})
	}
}

func TestTrackActivity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		pathID         string
		reqBody        string
		expectedTrack  domain.Track
		mockErr        error
		wantSrvCall    bool
		expectedStatus int
		expectedInBody string
	}{
		{
			name:           "valid activity track -> 200 OK",
			pathID:         validUUID,
			reqBody:        `{"date": "2026-11-16", "value": 45}`,
			expectedTrack:  domain.Track{Date: "2026-11-16", Value: 45},
			mockErr:        nil,
			wantSrvCall:    true,
			expectedStatus: http.StatusOK,
			expectedInBody: `{"status":"ok"}`,
		},
		{
			name:           "invalid UUID -> 400 Bad Request",
			pathID:         "123",
			reqBody:        `{"date": "2026-11-16", "value": 45}`,
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid ID"`,
		},
		{
			name:           "invalid JSON -> 400 Bad Request",
			pathID:         validUUID,
			reqBody:        `{"date": 12345}`,
			wantSrvCall:    false,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid JSON"`,
		},
		{
			name:           "service returns BadDate -> 400 Bad Request",
			pathID:         validUUID,
			reqBody:        `{"date": "2020-01-01", "value": 45}`,
			expectedTrack:  domain.Track{Date: "2020-01-01", Value: 45},
			mockErr:        coreErrors.BadDate(),
			wantSrvCall:    true,
			expectedStatus: http.StatusBadRequest,
			expectedInBody: `"invalid date"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srvCalled := false
			srv := &mockService{
				SaveDataFunc: func(ctx context.Context, id string, userData domain.Track) error {
					srvCalled = true
					assert.Equal(t, tc.pathID, id)
					assert.Equal(t, tc.expectedTrack, userData)
					return tc.mockErr
				},
			}

			h := transport.NewTransportHTTP(testLogger(t), srv)

			req := httptest.NewRequest(http.MethodPost, "/api/habits/"+tc.pathID+"/activities", bytes.NewBufferString(tc.reqBody))
			req.SetPathValue("id", tc.pathID)
			rec := httptest.NewRecorder()

			h.TrackActivity(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.wantSrvCall, srvCalled)
			assert.Contains(t, rec.Body.String(), tc.expectedInBody)
		})
	}
}
