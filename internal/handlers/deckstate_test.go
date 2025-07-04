package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/config"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
	"github.com/jwebster45206/tcg-api/internal/state"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestDeckStateHandler_CreateDeckState(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})

	// Create sample deck
	deckID := uuid.New()
	sampleDeck := &deckdef.Deck{
		ID:        deckID,
		Name:      "Test Deck",
		Type:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Cards: &deckdef.CardCollection{
			TotalCount:  3,
			UniqueCount: 2,
			Items: []deckdef.CardWithQuantity{
				{
					Card: &deckdef.ImageCard{
						ID:            uuid.New(),
						Name:          "Test Card 1",
						Description:   "Test description",
						FrontImageURL: "https://example.com/front.jpg",
						BackImageURL:  "https://example.com/back.jpg",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
					Quantity: 2,
				},
				{
					Card: &deckdef.ImageCard{
						ID:            uuid.New(),
						Name:          "Test Card 2",
						Description:   "Test description 2",
						FrontImageURL: "https://example.com/front2.jpg",
						BackImageURL:  "https://example.com/back2.jpg",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
					Quantity: 1,
				},
			},
		},
	}

	tests := []struct {
		name           string
		requestBody    CreateDeckStateRequest
		setupMocks     func(storage.Storage, *state.MockDeckStateStorage)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful creation",
			requestBody: CreateDeckStateRequest{
				DeckID:      deckID,
				PlayerCount: 2,
			},
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// Create the deck first
				_, err := mockStorage.CreateDeck(context.Background(), *sampleDeck)
				if err != nil {
					t.Errorf("Failed to create test deck: %v", err)
				}
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "deck not found",
			requestBody: CreateDeckStateRequest{
				DeckID:      uuid.New(),
				PlayerCount: 2,
			},
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// Don't create any deck - will return not found
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name: "invalid player count - too low",
			requestBody: CreateDeckStateRequest{
				DeckID:      deckID,
				PlayerCount: 0,
			},
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				_, err := mockStorage.CreateDeck(context.Background(), *sampleDeck)
				if err != nil {
					t.Errorf("Failed to create test deck: %v", err)
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid player count - too high",
			requestBody: CreateDeckStateRequest{
				DeckID:      deckID,
				PlayerCount: 9,
			},
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				_, err := mockStorage.CreateDeck(context.Background(), *sampleDeck)
				if err != nil {
					t.Errorf("Failed to create test deck: %v", err)
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing deck ID",
			requestBody: CreateDeckStateRequest{
				DeckID:      uuid.Nil,
				PlayerCount: 2,
			},
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// No setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "storage save error",
			requestBody: CreateDeckStateRequest{
				DeckID:      deckID,
				PlayerCount: 2,
			},
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				_, err := mockStorage.CreateDeck(context.Background(), *sampleDeck)
				if err != nil {
					t.Errorf("Failed to create test deck: %v", err)
				}
				mockStateStorage.SetSaveError(errors.New("redis connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockStorage := storage.NewMockStorage()
			mockStateStorage := state.NewMockDeckStateStorage()
			tt.setupMocks(mockStorage, mockStateStorage)

			// Create handler
			handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

			// Create request
			requestBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/deckstates", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check if error response format is correct
			if tt.expectError {
				var errorResponse ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse.Error == "" {
					t.Errorf("Expected error field to be populated")
				}
			} else {
				// Check if deck state was created successfully
				var deckState deckstate.DeckState
				if err := json.Unmarshal(rr.Body.Bytes(), &deckState); err != nil {
					t.Errorf("Failed to unmarshal deck state response: %v", err)
				}
				if deckState.ID == "" {
					t.Errorf("Expected deck state ID to be populated")
				}
				if deckState.PlayerCount != tt.requestBody.PlayerCount {
					t.Errorf("Expected player count %d, got %d", tt.requestBody.PlayerCount, deckState.PlayerCount)
				}
			}
		})
	}
}

func TestDeckStateHandler_GetDeckState(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})

	// Create sample deck state
	deckStateID := uuid.New().String()
	sampleDeckState := &deckstate.DeckState{
		ID:          deckStateID,
		PlayerCount: 2,
		Deck: deckdef.Deck{
			ID:   uuid.New(),
			Name: "Test Deck",
			Type: "standard",
		},
		Zones: map[string]deckstate.Zone{
			"draw": {
				Name:          "draw",
				Type:          deckstate.ZoneTypeDraw,
				DefaultFacing: deckstate.FaceDown,
				Items:         []deckstate.ZoneItem{},
			},
		},
	}

	tests := []struct {
		name           string
		deckStateID    string
		setupMocks     func(storage.Storage, *state.MockDeckStateStorage)
		expectedStatus int
		expectError    bool
	}{
		{
			name:        "successful get",
			deckStateID: deckStateID,
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				_ = mockStateStorage.SaveDeckState(context.Background(), deckStateID, sampleDeckState)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:        "deck state not found",
			deckStateID: uuid.New().String(),
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// Don't save any deck state
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:        "invalid UUID format",
			deckStateID: "invalid-uuid",
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// No setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:        "storage get error",
			deckStateID: deckStateID,
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				mockStateStorage.SetGetError(errors.New("redis connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockStorage := storage.NewMockStorage()
			mockStateStorage := state.NewMockDeckStateStorage()
			tt.setupMocks(mockStorage, mockStateStorage)

			// Create handler
			handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/v1/deckstates/"+tt.deckStateID, nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check response format
			if tt.expectError {
				var errorResponse ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse.Error == "" {
					t.Errorf("Expected error field to be populated")
				}
			} else {
				var deckState deckstate.DeckState
				if err := json.Unmarshal(rr.Body.Bytes(), &deckState); err != nil {
					t.Errorf("Failed to unmarshal deck state response: %v", err)
				}
				if deckState.ID != deckStateID {
					t.Errorf("Expected deck state ID %s, got %s", deckStateID, deckState.ID)
				}
			}
		})
	}
}

func TestDeckStateHandler_DeleteDeckState(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})

	// Create sample deck state
	deckStateID := uuid.New().String()
	sampleDeckState := &deckstate.DeckState{
		ID:          deckStateID,
		PlayerCount: 2,
		Deck: deckdef.Deck{
			ID:   uuid.New(),
			Name: "Test Deck",
			Type: "standard",
		},
		Zones: map[string]deckstate.Zone{
			"draw": {
				Name:          "draw",
				Type:          deckstate.ZoneTypeDraw,
				DefaultFacing: deckstate.FaceDown,
				Items:         []deckstate.ZoneItem{},
			},
		},
	}

	tests := []struct {
		name           string
		deckStateID    string
		setupMocks     func(storage.Storage, *state.MockDeckStateStorage)
		expectedStatus int
		expectError    bool
	}{
		{
			name:        "successful delete",
			deckStateID: deckStateID,
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				_ = mockStateStorage.SaveDeckState(context.Background(), deckStateID, sampleDeckState)
			},
			expectedStatus: http.StatusNoContent,
			expectError:    false,
		},
		{
			name:        "deck state not found",
			deckStateID: uuid.New().String(),
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// Don't save any deck state
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:        "invalid UUID format",
			deckStateID: "invalid-uuid",
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				// No setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:        "storage get error when checking existence",
			deckStateID: deckStateID,
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				mockStateStorage.SetGetError(errors.New("redis connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:        "storage delete error",
			deckStateID: deckStateID,
			setupMocks: func(mockStorage storage.Storage, mockStateStorage *state.MockDeckStateStorage) {
				_ = mockStateStorage.SaveDeckState(context.Background(), deckStateID, sampleDeckState)
				mockStateStorage.SetDeleteError(errors.New("redis delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockStorage := storage.NewMockStorage()
			mockStateStorage := state.NewMockDeckStateStorage()
			tt.setupMocks(mockStorage, mockStateStorage)

			// Create handler
			handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/v1/deckstates/"+tt.deckStateID, nil)

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check response format
			if tt.expectError {
				var errorResponse ErrorResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse.Error == "" {
					t.Errorf("Expected error field to be populated")
				}
			} else {
				// For successful delete, response should be empty with 204 status
				if rr.Body.Len() > 0 {
					t.Errorf("Expected empty response body for successful delete, got: %s", rr.Body.String())
				}

				// Verify the deck state was actually deleted
				if mockStateStorage.HasState(deckStateID) {
					t.Errorf("Expected deck state to be deleted, but it still exists")
				}
			}
		})
	}
}

func TestDeckStateHandler_InvalidMethods(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "PATCH method not allowed",
			method:         http.MethodPatch,
			path:           "/v1/deckstates",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT method not allowed",
			method:         http.MethodPut,
			path:           "/v1/deckstates",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "POST with ID not allowed",
			method:         http.MethodPost,
			path:           "/v1/deckstates/" + uuid.New().String(),
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "GET without ID not allowed",
			method:         http.MethodGet,
			path:           "/v1/deckstates",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "DELETE without ID not allowed",
			method:         http.MethodDelete,
			path:           "/v1/deckstates",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestDeckStateHandler_InvalidJSON(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	// Test with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/v1/deckstates", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var errorResponse ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponse); err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}
	if errorResponse.Error != errStrBadRequest {
		t.Errorf("Expected error '%s', got '%s'", errStrBadRequest, errorResponse.Error)
	}
}

func TestDeckStateHandler_AddZoneAction(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	deckStateID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/v1/deckstates/"+deckStateID+"/actions/add-zone", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expected status %d, got %d", http.StatusNotImplemented, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Error != errStrNotImplemented {
		t.Errorf("Expected error '%s', got '%s'", errStrNotImplemented, response.Error)
	}
}

func TestDeckStateHandler_RemoveZoneAction(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	deckStateID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/v1/deckstates/"+deckStateID+"/actions/remove-zone", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expected status %d, got %d", http.StatusNotImplemented, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Error != errStrNotImplemented {
		t.Errorf("Expected error '%s', got '%s'", errStrNotImplemented, response.Error)
	}
}

func TestDeckStateHandler_SortZoneAction(t *testing.T) {
	logger := config.NewLogger(config.LoggerConfig{Level: "ERROR"})
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	deckStateID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/v1/deckstates/"+deckStateID+"/actions/sort-zone", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expected status %d, got %d", http.StatusNotImplemented, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Error != errStrNotImplemented {
		t.Errorf("Expected error '%s', got '%s'", errStrNotImplemented, response.Error)
	}
}
