package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/N1ktarchik/habit-graph/internal/core/transport/middleware"
)

func TestCorsMiddleware(t *testing.T) {
	t.Parallel()

	const allowedOrigin = "http://localhost:3000"

	testCases := []struct {
		name           string
		method         string
		wantNextCalled bool
		expectedStatus int
	}{
		{
			name:           "preflight OPTIONS request returns 204 and skips next handler",
			method:         http.MethodOptions,
			wantNextCalled: false,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "GET request sets CORS headers and calls next handler",
			method:         http.MethodGet,
			wantNextCalled: true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request sets CORS headers and calls next handler",
			method:         http.MethodPost,
			wantNextCalled: true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := middleware.CorsMiddleware(nextHandler, allowedOrigin)

			req := httptest.NewRequest(tc.method, "/api/habits", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.wantNextCalled, nextCalled)
			assert.Equal(t, allowedOrigin, rec.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET, POST, PATCH, DELETE, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
		})
	}
}
