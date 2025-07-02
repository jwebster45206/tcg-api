package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jwebster45206/tcg-api/internal/handlers"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func createTestRedisClient() *redis.Client {
	// Start an in-memory Redis server for testing
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	// Create a Redis client that connects to the mock server
	return redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
}

func TestMainRoutes(t *testing.T) {
	sto := storage.NewMockStorage()
	redisClient := createTestRedisClient()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.NewHealthHandler(sto, redisClient))

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
	sto := storage.NewMockStorage()

	mockRedis := createTestRedisClient()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.NewHealthHandler(sto, mockRedis))

	server := &http.Server{
		Addr:         ":0", // Use port 0 to get any available port
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Test that the server has the expected configuration
	if server.Handler != mux {
		t.Error("server handler should match the mux we created")
	}

	if server.Addr != ":0" {
		t.Error("server address should be :0")
	}
}
