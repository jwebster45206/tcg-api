package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// createTestRedisClient creates a Redis client that connects to a mock Redis server
// Using miniredis is acceptable for unit tests when testing Redis integration
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

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	sto := storage.NewMockStorage()
	redisClient := createTestRedisClient()
	handler := NewHealthHandler(sto, redisClient)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHealthHandler_WithRedis(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	sto := storage.NewMockStorage()
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
