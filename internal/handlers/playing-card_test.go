package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestPlayingCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/playing-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

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

func TestPlayingCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New().String()
	req, err := http.NewRequest("GET", "/playing-cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

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

func TestPlayingCardsHandler_GetCard_InvalidID(t *testing.T) {
	// Test with invalid UUID
	req, err := http.NewRequest("GET", "/playing-cards/invalid-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

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

func TestPlayingCardsHandler_CreateCard(t *testing.T) {
	cardReq := models.PlayingCard{
		Suite:         models.SuiteHearts,
		Ranking:       1, // Ace
		FrontImageURL: "https://example.com/ace_hearts_front.png",
		BackImageURL:  "https://example.com/card_back.png",
	}

	jsonBody, _ := json.Marshal(cardReq)
	req, err := http.NewRequest("POST", "/playing-cards", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdCard models.PlayingCard
	if err := json.Unmarshal(rr.Body.Bytes(), &createdCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if createdCard.Suite != cardReq.Suite {
		t.Errorf("Expected card suite '%s', got '%s'", cardReq.Suite, createdCard.Suite)
	}

	if createdCard.Ranking != cardReq.Ranking {
		t.Errorf("Expected card ranking %d, got %d", cardReq.Ranking, createdCard.Ranking)
	}

	if createdCard.ID == uuid.Nil {
		t.Error("Expected card to have a generated ID")
	}
}

func TestPlayingCardsHandler_CreateCard_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/playing-cards", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

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

func TestPlayingCardsHandler_CreateCard_InvalidCard(t *testing.T) {
	// Test with invalid ranking (out of range)
	cardReq := models.PlayingCard{
		Suite:         models.SuiteHearts,
		Ranking:       15, // Invalid - should be 1-13
		FrontImageURL: "https://example.com/front.png",
		BackImageURL:  "https://example.com/back.png",
	}

	jsonBody, _ := json.Marshal(cardReq)
	req, err := http.NewRequest("POST", "/playing-cards", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "invalid_card" {
		t.Errorf("Expected error 'invalid_card', got '%s'", response.Error)
	}
}

func TestPlayingCardsHandler_UpdateCard(t *testing.T) {
	cardReq := models.PlayingCard{
		Suite:   models.SuiteSpades,
		Ranking: 10,
	}

	mockStorage := storage.NewMockStorage()
	logger := testLogger()

	// Create the card first
	ctx := context.Background()
	cardReq.ID = uuid.New() // Ensure we have a valid ID for the test
	_, err := mockStorage.CreatePlayingCard(ctx, cardReq)
	if err != nil {
		t.Fatalf("Failed to create test card: %v", err)
	}

	// Now update it
	updateReq := models.PlayingCard{
		Suite:         models.SuiteDiamonds,
		Ranking:       5,
		FrontImageURL: "https://example.com/updated_front.png",
		BackImageURL:  "https://example.com/updated_back.png",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", "/playing-cards/"+cardReq.ID.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := NewPlayingCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedCard models.PlayingCard
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if updatedCard.Suite != updateReq.Suite {
		t.Errorf("Expected updated card suite '%s', got '%s'", updateReq.Suite, updatedCard.Suite)
	}

	if updatedCard.Ranking != updateReq.Ranking {
		t.Errorf("Expected updated card ranking %d, got %d", updateReq.Ranking, updatedCard.Ranking)
	}
}

func TestPlayingCardsHandler_UpdateCard_NotFound(t *testing.T) {
	cardID := uuid.New().String()
	updateReq := models.PlayingCard{
		Suite:   models.SuiteClubs,
		Ranking: 7,
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", "/playing-cards/"+cardID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

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

func TestPlayingCardsHandler_UpdateCard_InvalidCard(t *testing.T) {
	cardID := uuid.New().String()
	updateReq := models.PlayingCard{
		Suite:   "invalid_suite", // Invalid suite
		Ranking: 5,
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", "/playing-cards/"+cardID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "invalid_card" {
		t.Errorf("Expected error 'invalid_card', got '%s'", response.Error)
	}
}

func TestPlayingCardsHandler_DeleteCard(t *testing.T) {
	cardReq := models.PlayingCard{
		Suite:   models.SuiteHearts,
		Ranking: 12, // Queen
	}

	mockStorage := storage.NewMockStorage()
	logger := testLogger()

	// Create the card first
	ctx := context.Background()
	cardReq.ID = uuid.New() // Ensure we have a valid ID for the test
	_, err := mockStorage.CreatePlayingCard(ctx, cardReq)
	if err != nil {
		t.Fatalf("Failed to create test card: %v", err)
	}

	req, err := http.NewRequest("DELETE", "/playing-cards/"+cardReq.ID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := NewPlayingCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestPlayingCardsHandler_DeleteCard_NotFound(t *testing.T) {
	cardID := uuid.New().String()
	req, err := http.NewRequest("DELETE", "/playing-cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

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

func TestPlayingCardsHandler_UnsupportedMethod(t *testing.T) {
	req, err := http.NewRequest("PATCH", "/playing-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewPlayingCardsHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
