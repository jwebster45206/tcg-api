package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	sto := storage.NewMockStorage()

	// Create a mock Redis client (nil for testing purposes)
	// In a real test environment, you might want to use a Redis mock library
	// or testcontainers, but for now we'll test with nil
	var redisClient *redis.Client = nil

	handler := NewHealthHandler(sto, redisClient)

	handler.ServeHTTP(rr, req)

	// With nil Redis client, the service should be unhealthy
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusServiceUnavailable)
	}
}

func TestHealthHandler_WithRedis(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	sto := storage.NewMockStorage()

	// Create a Redis client that points to a non-existent server
	// This will test the Redis connection failure case
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:9999", // Non-existent Redis server
	})

	handler := NewHealthHandler(sto, redisClient)

	handler.ServeHTTP(rr, req)

	// With failing Redis connection, the service should be unhealthy
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusServiceUnavailable)
	}
}
