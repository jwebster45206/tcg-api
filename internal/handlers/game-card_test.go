package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestGameCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/game-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
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

	// Should return an empty list since mock storage starts empty
	if len(cards) != 0 {
		t.Errorf("Expected empty card list, got %d cards", len(cards))
	}
}

func TestGameCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New().String()
	req, err := http.NewRequest("GET", "/game-cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
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
	req, err := http.NewRequest("GET", "/game-cards/invalid-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
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

	if response.Error != "invalid_id" {
		t.Errorf("Expected error 'invalid_id', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_CreateCard(t *testing.T) {
	cardReq := models.GameCard{
		Name:       "Test Card",
		Subtitle:   "A test card",
		Cost:       3,
		Type:       "Creature",
		Offense:    2,
		Defense:    3,
		Keywords:   []string{"Flying"},
		Colors:     []string{"Blue"},
		IsResource: false,
	}

	jsonBody, _ := json.Marshal(cardReq)
	req, err := http.NewRequest("POST", "/game-cards", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdCard models.GameCard
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
	req, err := http.NewRequest("POST", "/game-cards", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
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

	if response.Error != "invalid_json" {
		t.Errorf("Expected error 'invalid_json', got '%s'", response.Error)
	}
}

func TestGameCardsHandler_UpdateCard(t *testing.T) {
	cardReq := models.GameCard{
		Name: "Original Card",
		Type: "Creature",
	}

	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Create the card first
	ctx := context.Background()
	cardReq.ID = uuid.New() // Ensure we have a valid ID for the test
	_, err := mockStorage.CreateGameCard(ctx, cardReq)
	if err != nil {
		t.Fatalf("Failed to create test card: %v", err)
	}

	// Now update it
	updateReq := models.GameCard{
		Name: "Updated Card",
		Type: "Instant",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", "/game-cards/"+cardReq.ID.String(), bytes.NewBuffer(jsonBody))
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

	var updatedCard models.GameCard
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if updatedCard.Name != updateReq.Name {
		t.Errorf("Expected updated card name '%s', got '%s'", updateReq.Name, updatedCard.Name)
	}
}

func TestGameCardsHandler_UpdateCard_NotFound(t *testing.T) {
	cardID := uuid.New().String()
	updateReq := models.GameCard{
		Name: "Updated Card",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", "/game-cards/"+cardID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
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
	cardReq := models.GameCard{
		Name: "Card to Delete",
		Type: "Creature",
	}

	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Create the card first
	ctx := context.Background()
	cardReq.ID = uuid.New() // Ensure we have a valid ID for the test
	_, err := mockStorage.CreateGameCard(ctx, cardReq)
	if err != nil {
		t.Fatalf("Failed to create test card: %v", err)
	}

	req, err := http.NewRequest("DELETE", "/game-cards/"+cardReq.ID.String(), nil)
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
	req, err := http.NewRequest("DELETE", "/game-cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
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
	req, err := http.NewRequest("PATCH", "/game-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	handler := NewGameCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
