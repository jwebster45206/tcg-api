package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	sto := storage.NewMockStorage()
	handler := NewHealthHandler(sto)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
