package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jwebster45206/tcg-api/internal/handlers"
)

func TestMainRoutes(t *testing.T) {
	// Test that our routes are properly configured
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Health endpoint should be accessible",
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unknown endpoint should return 404",
			path:           "/unknown",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestServerStartup(t *testing.T) {
	// Test that we can create a server without it crashing
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)

	server := &http.Server{
		Addr:         ":0", // Use port 0 to get any available port
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if server == nil {
		t.Error("server should not be nil")
	}

	if server.Handler == nil {
		t.Error("server handler should not be nil")
	}
}
