package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/storage"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
)

func TestPlayingCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/playing-cards", nil)
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

	// Should return 3 cards since mock storage has sample data
	if len(cards) != 3 {
		t.Errorf("Expected 3 cards in list, got %d cards", len(cards))
	}
}

func TestPlayingCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New().String()
	req, err := http.NewRequest("GET", "/v1/playing-cards/"+cardID, nil)
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
	req, err := http.NewRequest("GET", "/v1/playing-cards/invalid-id", nil)
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

	if response.Error != "bad_request" {
		t.Errorf("Expected error 'bad_request', got '%s'", response.Error)
	}
}

func TestPlayingCardsHandler_CreateCard(t *testing.T) {
	cardReq := deckdef.PlayingCard{
		Suit:          deckdef.SuitHearts,
		Ranking:       1, // Ace
		FrontImageURL: "https://example.com/ace_hearts_front.png",
		BackImageURL:  "https://example.com/card_back.png",
	}

	jsonBody, _ := json.Marshal(cardReq)
	req, err := http.NewRequest("POST", "/v1/playing-cards", bytes.NewBuffer(jsonBody))
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

	var createdCard deckdef.PlayingCard
	if err := json.Unmarshal(rr.Body.Bytes(), &createdCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if createdCard.Suit != cardReq.Suit {
		t.Errorf("Expected card suit '%s', got '%s'", cardReq.Suit, createdCard.Suit)
	}

	if createdCard.Ranking != cardReq.Ranking {
		t.Errorf("Expected card ranking %d, got %d", cardReq.Ranking, createdCard.Ranking)
	}

	if createdCard.ID == uuid.Nil {
		t.Error("Expected card to have a generated ID")
	}
}

func TestPlayingCardsHandler_CreateCard_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/playing-cards", bytes.NewBuffer([]byte("invalid json")))
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

	if response.Error != "bad_request" {
		t.Errorf("Expected error 'bad_request', got '%s'", response.Error)
	}
}

func TestPlayingCardsHandler_CreateCard_InvalidCard(t *testing.T) {
	// Test with invalid ranking (out of range)
	cardReq := deckdef.PlayingCard{
		Suit:          deckdef.SuitHearts,
		Ranking:       15, // Invalid - should be 1-13
		FrontImageURL: "https://example.com/front.png",
		BackImageURL:  "https://example.com/back.png",
	}

	jsonBody, _ := json.Marshal(cardReq)
	req, err := http.NewRequest("POST", "/v1/playing-cards", bytes.NewBuffer(jsonBody))
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

	if response.Error != "validation_error" {
		t.Errorf("Expected error 'validation_error', got '%s'", response.Error)
	}
}

func TestPlayingCardsHandler_UpdateCard(t *testing.T) {
	cardReq := deckdef.PlayingCard{
		Suit:    deckdef.SuitSpades,
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
	updateReq := deckdef.PlayingCard{
		Suit:          deckdef.SuitDiamonds,
		Ranking:       5,
		FrontImageURL: "https://example.com/updated_front.png",
		BackImageURL:  "https://example.com/updated_back.png",
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PATCH", "/v1/playing-cards/"+cardReq.ID.String(), bytes.NewBuffer(jsonBody))
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

	var updatedCard deckdef.PlayingCard
	if err := json.Unmarshal(rr.Body.Bytes(), &updatedCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if updatedCard.Suit != updateReq.Suit {
		t.Errorf("Expected updated card suit '%s', got '%s'", updateReq.Suit, updatedCard.Suit)
	}

	if updatedCard.Ranking != updateReq.Ranking {
		t.Errorf("Expected updated card ranking %d, got %d", updateReq.Ranking, updatedCard.Ranking)
	}
}

func TestPlayingCardsHandler_UpdateCard_NotFound(t *testing.T) {
	cardID := uuid.New().String()
	updateReq := deckdef.PlayingCard{
		Suit:    deckdef.SuitClubs,
		Ranking: 7,
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PATCH", "/v1/playing-cards/"+cardID, bytes.NewBuffer(jsonBody))
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

func TestPlayingCardsHandler_UpdateCard_InvalidCard(t *testing.T) {
	// First create a valid card
	cardReq := deckdef.PlayingCard{
		Suit:    deckdef.SuitHearts,
		Ranking: 5,
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

	// Now try to update it with invalid data
	updateReq := deckdef.PlayingCard{
		ID:      cardReq.ID,     // Use the same ID as the created card
		Suit:    "invalid_suit", // Invalid suit
		Ranking: 5,
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PATCH", "/v1/playing-cards/"+cardReq.ID.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
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

	if response.Error != "validation_error" {
		t.Errorf("Expected error 'validation_error', got '%s'", response.Error)
	}
}

func TestPlayingCardsHandler_DeleteCard(t *testing.T) {
	cardReq := deckdef.PlayingCard{
		Suit:    deckdef.SuitHearts,
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

	req, err := http.NewRequest("DELETE", "/v1/playing-cards/"+cardReq.ID.String(), nil)
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
	req, err := http.NewRequest("DELETE", "/v1/playing-cards/"+cardID, nil)
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
	req, err := http.NewRequest("PUT", "/v1/playing-cards", nil)
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
