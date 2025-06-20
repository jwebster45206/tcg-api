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

func TestImageCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/image-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create handler with dependencies
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

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

func TestImageCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New()

	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	// Test getting non-existent card
	req, err := http.NewRequest("GET", "/image-cards/"+cardID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Test with invalid UUID format
	req, err = http.NewRequest("GET", "/image-cards/invalid-uuid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestImageCardsHandler_CreateCard(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	card := models.ImageCard{
		Name:          "Test Image Card",
		Description:   "A test image card",
		FrontImageURL: "https://example.com/front.jpg",
		BackImageURL:  "https://example.com/back.jpg",
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/image-cards", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var createdCard models.ImageCard
	if err := json.Unmarshal(rr.Body.Bytes(), &createdCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if createdCard.Name != card.Name {
		t.Errorf("Expected card name %s, got %s", card.Name, createdCard.Name)
	}

	if createdCard.ID == uuid.Nil {
		t.Error("Expected created card to have an ID")
	}
}

func TestImageCardsHandler_CreateCard_InvalidJSON(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	req, err := http.NewRequest("POST", "/image-cards", bytes.NewReader([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestImageCardsHandler_UpdateCard(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	// First create a card
	card := models.ImageCard{
		Name:          "Original Image Card",
		Description:   "Original description",
		FrontImageURL: "https://example.com/front.jpg",
		BackImageURL:  "https://example.com/back.jpg",
	}

	createdCard, err := mockStorage.CreateImageCard(context.Background(), card)
	if err != nil {
		t.Fatal(err)
	}

	// Update the card
	updatedCard := models.ImageCard{
		ID:            createdCard.ID,
		Name:          "Updated Image Card",
		Description:   "Updated description",
		FrontImageURL: "https://example.com/new-front.jpg",
		BackImageURL:  "https://example.com/new-back.jpg",
	}

	jsonData, err := json.Marshal(updatedCard)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/image-cards/"+createdCard.ID.String(), bytes.NewReader(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseCard models.ImageCard
	if err := json.Unmarshal(rr.Body.Bytes(), &responseCard); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if responseCard.Name != updatedCard.Name {
		t.Errorf("Expected card name %s, got %s", updatedCard.Name, responseCard.Name)
	}
}

func TestImageCardsHandler_DeleteCard(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	// First create a card
	card := models.ImageCard{
		Name:          "Card to Delete",
		Description:   "This card will be deleted",
		FrontImageURL: "https://example.com/front.jpg",
		BackImageURL:  "https://example.com/back.jpg",
	}

	createdCard, err := mockStorage.CreateImageCard(context.Background(), card)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/image-cards/"+createdCard.ID.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the card was deleted
	_, err = mockStorage.GetImageCard(context.Background(), createdCard.ID)
	if err == nil {
		t.Error("Expected card to be deleted, but it still exists")
	}
}

func TestImageCardsHandler_MethodNotAllowed(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	logger := testLogger()
	handler := NewImageCardsHandler(mockStorage, logger)

	req, err := http.NewRequest("PATCH", "/image-cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
