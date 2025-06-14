package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

func TestCardsHandler_ListCards(t *testing.T) {
	req, err := http.NewRequest("GET", "/cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "not_implemented" {
		t.Errorf("Expected error 'not_implemented', got '%s'", response.Error)
	}
}

func TestCardsHandler_GetCard(t *testing.T) {
	// Test with valid UUID
	cardID := uuid.New().String()
	req, err := http.NewRequest("GET", "/cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "not_implemented" {
		t.Errorf("Expected error 'not_implemented', got '%s'", response.Error)
	}
}

func TestCardsHandler_GetCard_InvalidID(t *testing.T) {
	// Test with invalid UUID
	req, err := http.NewRequest("GET", "/cards/invalid-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

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

func TestCardsHandler_CreateCard(t *testing.T) {
	cardReq := models.Card{
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
	req, err := http.NewRequest("POST", "/cards", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "not_implemented" {
		t.Errorf("Expected error 'not_implemented', got '%s'", response.Error)
	}
}

func TestCardsHandler_CreateCard_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/cards", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

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

func TestCardsHandler_UpdateCard(t *testing.T) {
	cardID := uuid.New().String()
	name := "Updated Card"
	updateReq := models.Card{
		Name: name,
	}

	jsonBody, _ := json.Marshal(updateReq)
	req, err := http.NewRequest("PUT", "/cards/"+cardID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "not_implemented" {
		t.Errorf("Expected error 'not_implemented', got '%s'", response.Error)
	}
}

func TestCardsHandler_DeleteCard(t *testing.T) {
	cardID := uuid.New().String()
	req, err := http.NewRequest("DELETE", "/cards/"+cardID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotImplemented {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotImplemented)
	}

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse response body: %v", err)
	}

	if response.Error != "not_implemented" {
		t.Errorf("Expected error 'not_implemented', got '%s'", response.Error)
	}
}

func TestCardsHandler_UnsupportedMethod(t *testing.T) {
	req, err := http.NewRequest("PATCH", "/cards", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CardsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}
