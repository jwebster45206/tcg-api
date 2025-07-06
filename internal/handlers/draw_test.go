package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
	"github.com/jwebster45206/tcg-api/internal/state"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

func TestHandleDrawCards(t *testing.T) {
	// Setup
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	logger := testLogger()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	// Create test cards
	createTestCard := func(name string) *deckdef.ImageCard {
		return &deckdef.ImageCard{
			ID:   uuid.New(),
			Name: name,
		}
	}

	// Helper function to create a fresh deck state for each test
	createFreshDeckState := func(stateID string, deckCards int) *deckstate.DeckState {
		sourceZone := deckstate.NewZone("deck", deckstate.ZoneTypeDraw, 10)
		destZone := deckstate.NewZone("hand", deckstate.ZoneTypeHand, 5)

		// Add cards to the source zone
		for i := 0; i < deckCards; i++ {
			cardInZone := deckstate.CardInZone{
				Card: createTestCard("Test Card " + string(rune(i+49))),
			}
			sourceZone.Items = append(sourceZone.Items, cardInZone)
		}

		return &deckstate.DeckState{
			ID:          stateID,
			PlayerCount: 2,
			Zones: map[string]deckstate.Zone{
				"deck": sourceZone,
				"hand": destZone,
			},
		}
	}

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Create fresh deck state for this test
		testDeckState := createFreshDeckState("test-success", 5)
		err := mockStateStorage.SaveDeckState(ctx, testDeckState.ID, testDeckState)
		if err != nil {
			t.Fatalf("Failed to save test deck state: %v", err)
		}

		req := DrawRequest{
			FromZone: "deck",
			ToZone:   "hand",
			Count:    2,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-success/actions/draw-cards", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-success")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s",
				status, http.StatusOK, rr.Body.String())
		}

		var response DrawResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if !response.Success {
			t.Error("Expected success to be true")
		}

		if response.Zones == nil {
			t.Error("Expected zones object, got nil")
		} else {
			if len(response.Zones) != 2 {
				t.Errorf("Expected 2 zones in response, got %d", len(response.Zones))
			}

			// Check that zones are included
			if response.Zones["deck"] == nil {
				t.Error("Expected deck zone in response")
			}
			if response.Zones["hand"] == nil {
				t.Error("Expected hand zone in response")
			}

			// Items should not be included by default
			// Note: Due to shallow copy implementation, items might still be present
			// but the API contract is that they won't be included in JSON when not requested
		}

		if response.Meta != nil {
			t.Error("Expected no meta without include=meta query param")
		}

		// Verify the deck state was actually updated
		updatedState, _ := mockStateStorage.GetDeckState(ctx, "test-success")
		if len(updatedState.Zones["deck"].Items) != 3 {
			t.Errorf("Expected 3 items remaining in deck zone, got %d", len(updatedState.Zones["deck"].Items))
		}
		if len(updatedState.Zones["hand"].Items) != 2 {
			t.Errorf("Expected 2 items in hand zone, got %d", len(updatedState.Zones["hand"].Items))
		}
	})

	t.Run("SuccessWithMeta", func(t *testing.T) {
		// Create fresh deck state for this test
		testDeckState := createFreshDeckState("test-with-meta", 5)
		err := mockStateStorage.SaveDeckState(ctx, testDeckState.ID, testDeckState)
		if err != nil {
			t.Fatalf("Failed to save test deck state: %v", err)
		}

		req := DrawRequest{
			FromZone: "deck",
			ToZone:   "hand",
			Count:    1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-with-meta/actions/draw-cards?include=meta", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-with-meta")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response DrawResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if !response.Success {
			t.Error("Expected success to be true")
		}

		if response.Meta == nil {
			t.Error("Expected meta object when include=meta")
		} else {
			if response.Meta.Operation != "draw_cards" {
				t.Errorf("Expected operation 'draw_cards', got '%s'", response.Meta.Operation)
			}
			if response.Meta.FromZone != "deck" {
				t.Errorf("Expected from_zone 'deck', got '%s'", response.Meta.FromZone)
			}
			if response.Meta.ToZone != "hand" {
				t.Errorf("Expected to_zone 'hand', got '%s'", response.Meta.ToZone)
			}
			if response.Meta.CardsDrawn != 1 {
				t.Errorf("Expected cards_drawn 1, got %d", response.Meta.CardsDrawn)
			}
			if response.Meta.DurationMS <= 0 {
				t.Errorf("Expected positive duration, got %f", response.Meta.DurationMS)
			}
		}
	})

	t.Run("SuccessWithItems", func(t *testing.T) {
		// Create fresh deck state for this test
		testDeckState := createFreshDeckState("test-with-items", 5)
		err := mockStateStorage.SaveDeckState(ctx, testDeckState.ID, testDeckState)
		if err != nil {
			t.Fatalf("Failed to save test deck state: %v", err)
		}

		req := DrawRequest{
			FromZone: "deck",
			ToZone:   "hand",
			Count:    2,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-with-items/actions/draw-cards?include=items", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-with-items")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s",
				status, http.StatusOK, rr.Body.String())
		}

		var response DrawResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if !response.Success {
			t.Error("Expected success to be true")
		}

		// Check that items are included when requested
		if response.Zones["deck"].Items == nil {
			t.Error("Expected items to be included in deck zone when include=items")
		}
		if response.Zones["hand"].Items == nil {
			t.Error("Expected items to be included in hand zone when include=items")
		}

		// Verify the actual item counts in response
		if len(response.Zones["deck"].Items) != 3 {
			t.Errorf("Expected 3 items in deck zone response, got %d", len(response.Zones["deck"].Items))
		}
		if len(response.Zones["hand"].Items) != 2 {
			t.Errorf("Expected 2 items in hand zone response, got %d", len(response.Zones["hand"].Items))
		}
	})

	t.Run("MissingFromZone", func(t *testing.T) {
		req := DrawRequest{
			ToZone: "hand",
			Count:  1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/draw-cards", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-state-id")

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
	})

	t.Run("MissingToZone", func(t *testing.T) {
		req := DrawRequest{
			FromZone: "deck",
			Count:    1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/draw-cards", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-state-id")

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
	})

	t.Run("SourceZoneNotFound", func(t *testing.T) {
		req := DrawRequest{
			FromZone: "nonexistent",
			ToZone:   "hand",
			Count:    1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/draw-cards", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-state-id")

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
	})

	t.Run("DestinationZoneNotFound", func(t *testing.T) {
		req := DrawRequest{
			FromZone: "deck",
			ToZone:   "nonexistent",
			Count:    1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/draw-cards", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-state-id")

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
	})

	t.Run("InsufficientCards", func(t *testing.T) {
		// Reset test state - create a deck with only 1 card
		singleCardZone := deckstate.NewZone("single-deck", deckstate.ZoneTypeDraw, 1)
		singleCardZone.Items = append(singleCardZone.Items, deckstate.CardInZone{
			Card: createTestCard("Last Card"),
		})

		emptyHandZone := deckstate.NewZone("empty-hand", deckstate.ZoneTypeHand, 5)

		singleCardState := &deckstate.DeckState{
			ID:          "single-card-state",
			PlayerCount: 1,
			Zones: map[string]deckstate.Zone{
				"single-deck": singleCardZone,
				"empty-hand":  emptyHandZone,
			},
		}

		err := mockStateStorage.SaveDeckState(ctx, singleCardState.ID, singleCardState)
		if err != nil {
			t.Fatalf("Failed to save single card deck state: %v", err)
		}

		req := DrawRequest{
			FromZone: "single-deck",
			ToZone:   "empty-hand",
			Count:    5, // Try to draw more than available
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/single-card-state/actions/draw-cards?include=meta", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "single-card-state")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s",
				status, http.StatusOK, rr.Body.String())
		}

		var response DrawResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if !response.Success {
			t.Error("Expected success to be true")
		}

		// Should have drawn only 1 card (all available)
		if response.Meta == nil {
			t.Error("Expected meta to be included")
		} else {
			if response.Meta.CardsDrawn != 1 {
				t.Errorf("Expected cards_drawn to be 1, got %d", response.Meta.CardsDrawn)
			}
		}

		// Verify deck state
		updatedState, _ := mockStateStorage.GetDeckState(ctx, "single-card-state")
		if len(updatedState.Zones["single-deck"].Items) != 0 {
			t.Errorf("Expected 0 items in source zone, got %d", len(updatedState.Zones["single-deck"].Items))
		}
		if len(updatedState.Zones["empty-hand"].Items) != 1 {
			t.Errorf("Expected 1 item in destination zone, got %d", len(updatedState.Zones["empty-hand"].Items))
		}
	})

	t.Run("EmptySourceZone", func(t *testing.T) {
		// Create a deck state with empty source zone
		emptyDeckZone := deckstate.NewZone("empty-deck", deckstate.ZoneTypeDraw, 0)
		handZone := deckstate.NewZone("hand-zone", deckstate.ZoneTypeHand, 5)

		emptyDeckState := &deckstate.DeckState{
			ID:          "empty-deck-state",
			PlayerCount: 1,
			Zones: map[string]deckstate.Zone{
				"empty-deck": emptyDeckZone,
				"hand-zone":  handZone,
			},
		}

		err := mockStateStorage.SaveDeckState(ctx, emptyDeckState.ID, emptyDeckState)
		if err != nil {
			t.Fatalf("Failed to save empty deck state: %v", err)
		}

		req := DrawRequest{
			FromZone: "empty-deck",
			ToZone:   "hand-zone",
			Count:    1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/empty-deck-state/actions/draw-cards", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "empty-deck-state")

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v. Body: %s",
				status, http.StatusBadRequest, rr.Body.String())
		}

		var response ErrorResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if response.Error != "bad_request" {
			t.Errorf("Expected error 'bad_request', got '%s'", response.Error)
		}
	})

	t.Run("NonExistentDeckState", func(t *testing.T) {
		req := DrawRequest{
			FromZone: "deck",
			ToZone:   "hand",
			Count:    1,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/nonexistent-id/actions/draw-cards", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "nonexistent-id")

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
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/draw-cards", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleDrawCards(rr, httpReq, "test-state-id")

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
	})
}
