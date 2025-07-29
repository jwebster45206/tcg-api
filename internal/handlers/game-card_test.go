package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/storage"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
)

// testLogger creates a test logger that outputs to stdout but suppresses output during tests
func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors in tests
	}))
}

func TestGameCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/game-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var cards []interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &cards); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	// Should return 1 card since mock storage has sample data
	if len(cards) != 1 {
		t.Errorf("Expected 1 card in list, got %d cards", len(cards))
	}
}

func TestGameCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New().String()
	req, err := http.NewRequest("GET", "/v1/game-cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "not_found" {
		t.Errorf("Expected error 'not_found', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_GetCard_InvalidID(t *testing.T) {
	// Test with invalid UUID
	req, err := http.NewRequest("GET", "/v1/game-cards/invalid-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "bad_request" {
		t.Errorf("Expected error 'bad_request', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_CreateCard(t *testing.T) {
	cardReq := deckdef.GameCard{
		Name:        "Test Card",
		Description: "A test card",
		ManaCost:    3,
		CardType:    "creature",
		Attack:      2,
		Health:      3,
		Keywords:    []string{"flying"},
		Colors:      []string{"blue"},
		Rarity:      "common",
		SetCode:     "TEST",
	}

	jsonBody, _ := json.Marshal(cardReq)
	req, err := http.NewRequest("POST", "/v1/game-cards", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdCard deckdef.GameCard
	if err := json.Unmarshal(rr.Body.Bytes(), &createdCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if createdCard.Name != cardReq.Name {
		t.Errorf("Expected card name '%s', got '%s'", cardReq.Name, createdCard.Name)
	}

	if createdCard.ID == uuid.Nil {
		t.Error("Expected card to have a generated ID")
	}
}

func TestGameCardsHandler_CreateCard_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/game-cards", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "bad_request" {
		t.Errorf("Expected error 'bad_request', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_UpdateCard(t *testing.T) {
	cardReq := deckdef.GameCard{
		Name:     "Original Card",
		CardType: "creature",
	}

	mockStorage := storage.NewMockStorage()
	logger := testLogger()

	// Create the card first
	ctx := context.Background()
	cardReq.ID = uuid.New() // Ensure we have a valid ID for the test
	_, err := mockStorage.CreateGameCard(ctx, cardReq)
	if err != nil {
		t.Fatalf("Failed to create test card: %v", err)
	}

	// Now update it
	updateReq := deckdef.GameCard{
		Name:     "Updated Card",
		CardType: "spell",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PATCH", "/v1/game-cards/"+cardReq.ID.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedCard deckdef.GameCard
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if updatedCard.Name != updateReq.Name {
		t.Errorf("Expected updated card name '%s', got '%s'", updateReq.Name, updatedCard.Name)
	}
}

func TestGameCardsHandler_UpdateCard_NotFound(t *testing.T) {
	cardID := uuid.New().String()
	updateReq := deckdef.GameCard{
		Name: "Updated Card",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PATCH", "/v1/game-cards/"+cardID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "internal_error" {
		t.Errorf("Expected error 'internal_error', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_DeleteCard(t *testing.T) {
	cardReq := deckdef.GameCard{
		Name:     "Card to Delete",
		CardType: "creature",
	}

	mockStorage := storage.NewMockStorage()
	logger := testLogger()

	// Create the card first
	ctx := context.Background()
	cardReq.ID = uuid.New() // Ensure we have a valid ID for the test
	_, err := mockStorage.CreateGameCard(ctx, cardReq)
	if err != nil {
		t.Fatalf("Failed to create test card: %v", err)
	}

	req, err := http.NewRequest("DELETE", "/v1/game-cards/"+cardReq.ID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestGameCardsHandler_DeleteCard_NotFound(t *testing.T) {
	cardID := uuid.New().String()
	req, err := http.NewRequest("DELETE", "/v1/game-cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "internal_error" {
		t.Errorf("Expected error 'internal_error', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_UnsupportedMethod(t *testing.T) {
	req, err := http.NewRequest("PUT", "/v1/game-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
