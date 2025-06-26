package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestDecksHandler_ListDecks(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/decks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var decks []interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &decks); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Should return empty array for new mock storage
	if len(decks) != 0 {
		t.Errorf("Expected empty array, got %d decks", len(decks))
	}
}

func TestDecksHandler_CreateDeck(t *testing.T) {
	deckName := "Test Playing Deck"
	deckTypeVal := "playing-card"

	deck := models.Deck{
		Name:           deckName,
		DeckType:       &deckTypeVal,
		SleeveImageURL: stringPtr("https://example.com/sleeve.jpg"),
	}

	body, err := json.Marshal(deck)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/decks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var responseDeck models.Deck
	if err := json.Unmarshal(rr.Body.Bytes(), &responseDeck); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Verify the deck was created with expected values
	if responseDeck.Name != deckName {
		t.Errorf("Expected deck name %s, got %s", deckName, responseDeck.Name)
	}

	if responseDeck.DeckType == nil || *responseDeck.DeckType != deckTypeVal {
		t.Errorf("Expected deck type %s, got %v", deckTypeVal, responseDeck.DeckType)
	}

	if responseDeck.ID == uuid.Nil {
		t.Error("Expected deck ID to be generated")
	}

	if responseDeck.SleeveImageURL == nil || *responseDeck.SleeveImageURL != "https://example.com/sleeve.jpg" {
		t.Errorf("Expected sleeve image URL, got %v", responseDeck.SleeveImageURL)
	}
}

func TestDecksHandler_GetDeck(t *testing.T) {
	// Create a deck first
	deckID := uuid.New()
	deckName := "Test Deck"
	deckTypeVal := "standard"

	deck := models.Deck{
		ID:        deckID,
		Name:      deckName,
		DeckType:  &deckTypeVal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup mock storage with deck
	mockStorage := storage.NewMockStorage()
	_, err := mockStorage.CreateDeck(context.Background(), deck)
	if err != nil {
		t.Fatal(err)
	}

	// Test getting the deck
	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/decks/%s", deckID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseDeck models.Deck
	if err := json.Unmarshal(rr.Body.Bytes(), &responseDeck); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if responseDeck.ID != deckID {
		t.Errorf("Expected deck ID %s, got %s", deckID, responseDeck.ID)
	}

	if responseDeck.Name != deckName {
		t.Errorf("Expected deck name %s, got %s", deckName, responseDeck.Name)
	}
}

func TestDecksHandler_GetDeck_NotFound(t *testing.T) {
	nonExistentID := uuid.New()

	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/decks/%s", nonExistentID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with empty mock storage
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDecksHandler_GetDeck_InvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/decks/invalid-uuid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestDecksHandler_UpdateDeck(t *testing.T) {
	// Create a deck first
	deckID := uuid.New()
	originalName := "Original Deck"
	deckTypeVal := "standard"

	originalDeck := models.Deck{
		ID:        deckID,
		Name:      originalName,
		DeckType:  &deckTypeVal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup mock storage with deck
	mockStorage := storage.NewMockStorage()
	_, err := mockStorage.CreateDeck(context.Background(), originalDeck)
	if err != nil {
		t.Fatal(err)
	}

	// Update the deck
	updatedName := "Updated Deck Name"
	updatedTypeVal := "tcg"
	updateDeck := models.Deck{
		Name:           updatedName,
		DeckType:       &updatedTypeVal,
		SleeveImageURL: stringPtr("https://example.com/new-sleeve.jpg"),
	}

	body, err := json.Marshal(updateDeck)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/decks/%s", deckID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseDeck models.Deck
	if err := json.Unmarshal(rr.Body.Bytes(), &responseDeck); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if responseDeck.ID != deckID {
		t.Errorf("Expected deck ID %s, got %s", deckID, responseDeck.ID)
	}

	if responseDeck.Name != updatedName {
		t.Errorf("Expected updated name %s, got %s", updatedName, responseDeck.Name)
	}

	if responseDeck.DeckType == nil || *responseDeck.DeckType != updatedTypeVal {
		t.Errorf("Expected updated deck type %s, got %v", updatedTypeVal, responseDeck.DeckType)
	}
}

func TestDecksHandler_UpdateDeck_NotFound(t *testing.T) {
	nonExistentID := uuid.New()
	updateDeck := models.Deck{
		Name: "Updated Name",
	}

	body, err := json.Marshal(updateDeck)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/decks/%s", nonExistentID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDecksHandler_DeleteDeck(t *testing.T) {
	// Create a deck first
	deckID := uuid.New()
	deckName := "Deck to Delete"
	deckTypeVal := "standard"

	deck := models.Deck{
		ID:        deckID,
		Name:      deckName,
		DeckType:  &deckTypeVal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup mock storage with deck
	mockStorage := storage.NewMockStorage()
	_, err := mockStorage.CreateDeck(context.Background(), deck)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the deck
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/v1/decks/%s", deckID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the deck is actually deleted
	_, err = mockStorage.GetDeck(context.Background(), deckID)
	if err == nil {
		t.Error("Expected deck to be deleted, but it still exists")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got %v", err)
	}
}

func TestDecksHandler_DeleteDeck_NotFound(t *testing.T) {
	nonExistentID := uuid.New()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/v1/decks/%s", nonExistentID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDecksHandler_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{"PATCH", "/v1/decks"},
		{"POST", "/v1/decks/123"},
		{"PUT", "/v1/decks"},
		{"DELETE", "/v1/decks"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.method, tt.path), func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			mockStorage := storage.NewMockStorage()
			logger := testLogger()
			handler := NewDecksHandler(mockStorage, logger)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestDecksHandler_InvalidRequestBody(t *testing.T) {
	invalidJSON := `{"name": "Test", "invalid_json": }`

	req, err := http.NewRequest("POST", "/v1/decks", strings.NewReader(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewDecksHandler(mockStorage, logger)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
