package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// HealthHandler handles the health check endpoint
func NewHealthHandler(sto storage.Storage, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		checks := make(map[string]string)
		status := "healthy"
		httpStatus := http.StatusOK

		// Test database connection
		if sto != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Try to perform a simple query to test database connectivity
			// We'll use a lightweight query that most storage implementations can handle
			if err := testStorageConnection(ctx, sto); err != nil {
				checks["database"] = "unhealthy: " + err.Error()
				status = "unhealthy"
				httpStatus = http.StatusServiceUnavailable
			} else {
				checks["database"] = "healthy"
			}
		} else {
			checks["database"] = "unhealthy: no storage configured"
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		}

		// Test Redis connection
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := redisClient.Ping(ctx).Result()
			if err != nil {
				checks["redis"] = "unhealthy: " + err.Error()
				status = "unhealthy"
				httpStatus = http.StatusServiceUnavailable
			} else {
				checks["redis"] = "healthy"
			}
		} else {
			checks["redis"] = "unhealthy: no redis configured"
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		}

		response := HealthResponse{
			Status:    status,
			Timestamp: time.Now().UTC(),
			Service:   "tcg-api",
			Version:   "0.x",
			Checks:    checks,
		}

		writeJSONResponse(w, httpStatus, response)
	}
}

func testStorageConnection(ctx context.Context, sto storage.Storage) error {
	return sto.Ping(ctx)
}
