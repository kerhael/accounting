package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeHealthRepo struct {
	Err error
}

func (f FakeHealthRepo) Check(ctx context.Context) error {
	return f.Err
}

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		repoErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "DB OK",
			repoErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"server":"ok", "db": "ok"}`,
		},
		{
			name:           "DB Down",
			repoErr:        errors.New("DB down"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   `{"server": "ok", "db":"ko"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo := FakeHealthRepo{Err: tt.repoErr}
			handler := NewHealthHandler(fakeRepo)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.expectedStatus)
			}

			// Check body
			var got map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			want := make(map[string]string)
			_ = json.Unmarshal([]byte(tt.expectedBody), &want)

			if got["status"] != want["status"] {
				t.Errorf("got body %v, want %v", got, want)
			}
		})
	}
}
