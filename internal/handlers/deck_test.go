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
	deckType := "playing-card"

	deck := models.Deck{
		Name:           deckName,
		Type:           deckType,
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

	if responseDeck.Type != deckType {
		t.Errorf("Expected deck type %s, got %s", deckType, responseDeck.Type)
	}

	if responseDeck.ID == uuid.Nil {
		t.Error("Expected deck ID to be generated")
	}

	if responseDeck.SleeveImageURL == nil || *responseDeck.SleeveImageURL != "https://example.com/sleeve.jpg" {
		t.Errorf("Expected sleeve image URL, got %v", responseDeck.SleeveImageURL)
	}
}

func TestDecksHandler_CreateDeckWithCards(t *testing.T) {
	deckName := "Test Deck with Cards"
	deckType := "playing-card"

	// Use known card IDs from mock storage
	cardID1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001") // Ace of Spades
	cardID2 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002") // King of Hearts

	deckInput := models.DeckInput{
		Name:           deckName,
		DeckType:       deckType,
		SleeveImageURL: stringPtr("https://example.com/sleeve.jpg"),
		Cards: &models.CardCollectionInput{
			Items: []models.CardInputWithQuantity{
				{Card: models.CardInput{ID: cardID1}, Quantity: 2},
				{Card: models.CardInput{ID: cardID2}, Quantity: 1},
			},
		},
	}

	body, err := json.Marshal(deckInput)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/decks?include=cards", bytes.NewBuffer(body))
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
		t.Logf("Response body: %s", rr.Body.String())
	}

	var responseDeck models.Deck
	if err := json.Unmarshal(rr.Body.Bytes(), &responseDeck); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Verify the deck was created with expected values
	if responseDeck.Name != deckName {
		t.Errorf("Expected deck name %s, got %s", deckName, responseDeck.Name)
	}

	if responseDeck.Type != deckType {
		t.Errorf("Expected deck type %s, got %s", deckType, responseDeck.Type)
	}

	if responseDeck.ID == uuid.Nil {
		t.Error("Expected deck ID to be generated")
	}

	// Verify cards are included in response
	if responseDeck.Cards == nil {
		t.Error("Expected cards to be included in response")
	} else {
		if responseDeck.Cards.UniqueCount != 2 {
			t.Errorf("Expected 2 unique cards, got %d", responseDeck.Cards.UniqueCount)
		}
		if responseDeck.Cards.TotalCount != 3 {
			t.Errorf("Expected total count of 3, got %d", responseDeck.Cards.TotalCount)
		}
	}
}

func TestDecksHandler_GetDeck(t *testing.T) {
	// Create a deck first
	deckID := uuid.New()
	deckName := "Test Deck"
	deckType := "standard"

	deck := models.Deck{
		ID:        deckID,
		Name:      deckName,
		Type:      deckType,
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
	deckType := "standard"

	originalDeck := models.Deck{
		ID:        deckID,
		Name:      originalName,
		Type:      deckType,
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
	updatedType := "tcg"
	updateDeck := models.Deck{
		Name:           updatedName,
		Type:           updatedType,
		SleeveImageURL: stringPtr("https://example.com/new-sleeve.jpg"),
	}

	body, err := json.Marshal(updateDeck)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/v1/decks/%s", deckID), bytes.NewBuffer(body))
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

	if responseDeck.Type != updatedType {
		t.Errorf("Expected updated deck type %s, got %s", updatedType, responseDeck.Type)
	}
}

func TestDecksHandler_UpdateDeckWithCards(t *testing.T) {
	// Create a deck first
	deckName := "Test Deck for Update"
	deckType := "playing-card"

	deck := models.Deck{
		Name:           deckName,
		Type:           deckType,
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
		t.Fatalf("failed to create deck: got %v want %v", status, http.StatusCreated)
	}

	var createdDeck models.Deck
	if err := json.Unmarshal(rr.Body.Bytes(), &createdDeck); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Now update the deck with cards
	cardID1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001") // Ace of Spades
	cardID2 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003") // Queen of Diamonds

	deckInput := models.DeckInput{
		Name:     "Updated Deck Name",
		DeckType: deckType,
		Cards: &models.CardCollectionInput{
			Items: []models.CardInputWithQuantity{
				{Card: models.CardInput{ID: cardID1}, Quantity: 1},
				{Card: models.CardInput{ID: cardID2}, Quantity: 3},
			},
		},
	}

	updateBody, err := json.Marshal(deckInput)
	if err != nil {
		t.Fatal(err)
	}

	updateReq, err := http.NewRequest("PATCH", "/v1/decks/"+createdDeck.ID.String()+"?include=cards", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatal(err)
	}
	updateReq.Header.Set("Content-Type", "application/json")

	updateRr := httptest.NewRecorder()
	handler.ServeHTTP(updateRr, updateReq)

	if status := updateRr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Logf("Response body: %s", updateRr.Body.String())
	}

	var updatedDeck models.Deck
	if err := json.Unmarshal(updateRr.Body.Bytes(), &updatedDeck); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Verify the deck was updated
	if updatedDeck.Name != "Updated Deck Name" {
		t.Errorf("Expected deck name 'Updated Deck Name', got %s", updatedDeck.Name)
	}

	// Verify cards are included in response
	if updatedDeck.Cards == nil {
		t.Error("Expected cards to be included in response")
	} else {
		if updatedDeck.Cards.UniqueCount != 2 {
			t.Errorf("Expected 2 unique cards, got %d", updatedDeck.Cards.UniqueCount)
		}
		if updatedDeck.Cards.TotalCount != 4 {
			t.Errorf("Expected total count of 4, got %d", updatedDeck.Cards.TotalCount)
		}
	}
}

func TestDecksHandler_DeleteDeck(t *testing.T) {
	// Create a deck first
	deckID := uuid.New()
	deckName := "Deck to Delete"
	deckType := "standard"

	deck := models.Deck{
		ID:        deckID,
		Name:      deckName,
		Type:      deckType,
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
		{"PUT", "/v1/decks"},
		{"POST", "/v1/decks/123"},
		{"PATCH", "/v1/decks"},
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
