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

func TestHandleAddZone(t *testing.T) {
	// Setup
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	logger := testLogger()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	// Create a test deck state
	testDeckState := &deckstate.DeckState{
		ID:          "test-state-id",
		PlayerCount: 2,
		Zones:       make(map[string]deckstate.Zone),
	}

	ctx := context.Background()
	err := mockStateStorage.SaveDeckState(ctx, testDeckState.ID, testDeckState)
	if err != nil {
		t.Fatalf("Failed to save test deck state: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		req := AddZoneRequest{
			Name: "custom-pile",
			Type: deckstate.ZoneTypeTable,
			Size: &[]int{10}[0], // Helper to get pointer to int
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response AddZoneResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if response.Zone.Name != "custom-pile" {
			t.Errorf("Expected zone name 'custom-pile', got '%s'", response.Zone.Name)
		}

		if response.Zone.Type != deckstate.ZoneTypeTable {
			t.Errorf("Expected zone type 'table', got '%s'", response.Zone.Type)
		}

		if response.Meta != nil {
			t.Error("Expected no meta without include=meta query param")
		}
	})

	t.Run("SuccessWithMeta", func(t *testing.T) {
		req := AddZoneRequest{
			Name: "custom-pile-with-meta",
			Type: deckstate.ZoneTypeTable,
			Size: &[]int{10}[0], // Helper to get pointer to int
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone?include=meta", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response AddZoneResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if response.Zone.Name != "custom-pile-with-meta" {
			t.Errorf("Expected zone name 'custom-pile-with-meta', got '%s'", response.Zone.Name)
		}

		if response.Zone.Type != deckstate.ZoneTypeTable {
			t.Errorf("Expected zone type 'table', got '%s'", response.Zone.Type)
		}

		if response.Meta == nil {
			t.Error("Expected meta object when include=meta")
		} else if response.Meta.DurationMS <= 0 {
			t.Errorf("Expected positive duration, got %f", response.Meta.DurationMS)
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		// First, add a zone
		req1 := AddZoneRequest{
			Name: "duplicate-zone",
			Type: deckstate.ZoneTypeHand,
		}
		jsonBody1, _ := json.Marshal(req1)
		httpReq1, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody1))
		httpReq1.Header.Set("Content-Type", "application/json")
		rr1 := httptest.NewRecorder()
		handler.handleAddZone(rr1, httpReq1, "test-state-id")

		// Try to add the same zone again
		req2 := AddZoneRequest{
			Name: "duplicate-zone",
			Type: deckstate.ZoneTypeDiscard,
		}
		jsonBody2, _ := json.Marshal(req2)
		httpReq2, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody2))
		httpReq2.Header.Set("Content-Type", "application/json")
		rr2 := httptest.NewRecorder()
		handler.handleAddZone(rr2, httpReq2, "test-state-id")

		if status := rr2.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rr2.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if response.Error != "bad_request" {
			t.Errorf("Expected error 'bad_request', got '%s'", response.Error)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		req := AddZoneRequest{
			Type: deckstate.ZoneTypeHand,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("InvalidSizeHint", func(t *testing.T) {
		req := AddZoneRequest{
			Name: "test-zone",
			Type: deckstate.ZoneTypeHand,
			Size: &[]int{-5}[0], // Negative size hint
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("SuccessWithDefaultFacing", func(t *testing.T) {
		facing := deckstate.FaceDown
		req := AddZoneRequest{
			Name:          "face-down-pile",
			Type:          deckstate.ZoneTypeTable,
			DefaultFacing: &facing,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response AddZoneResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if response.Zone.DefaultFacing != deckstate.FaceDown {
			t.Errorf("Expected default facing 'face-down', got '%s'", response.Zone.DefaultFacing)
		}
	})

	t.Run("SuccessMinimalRequest", func(t *testing.T) {
		req := AddZoneRequest{
			Name: "minimal-zone",
			Type: deckstate.ZoneTypeHand,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var response AddZoneResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if response.Zone.Name != "minimal-zone" {
			t.Errorf("Expected zone name 'minimal-zone', got '%s'", response.Zone.Name)
		}

		if response.Zone.Type != deckstate.ZoneTypeHand {
			t.Errorf("Expected zone type 'hand', got '%s'", response.Zone.Type)
		}

		// Check that default values are set correctly
		if response.Zone.DefaultFacing != deckstate.InHand {
			t.Errorf("Expected default facing 'in-hand' for hand zone, got '%s'", response.Zone.DefaultFacing)
		}
	})

	t.Run("MissingType", func(t *testing.T) {
		req := AddZoneRequest{
			Name: "test-zone",
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

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

	t.Run("InvalidJSON", func(t *testing.T) {
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/add-zone", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "test-state-id")

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

	t.Run("NonExistentDeckState", func(t *testing.T) {
		req := AddZoneRequest{
			Name: "test-zone",
			Type: deckstate.ZoneTypeHand,
		}

		jsonBody, _ := json.Marshal(req)
		httpReq, _ := http.NewRequest("POST", "/v1/deckstates/nonexistent-id/actions/add-zone", bytes.NewBuffer(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.handleAddZone(rr, httpReq, "nonexistent-id")

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
}

func TestHandleRemoveZone(t *testing.T) {
	// Setup
	mockStorage := storage.NewMockStorage()
	mockStateStorage := state.NewMockDeckStateStorage()
	logger := testLogger()
	handler := NewDeckStateHandler(mockStorage, mockStateStorage, logger)

	// Create a test deck state with zones
	testDeckState := &deckstate.DeckState{
		ID:          "test-state-id",
		PlayerCount: 2,
		Zones:       make(map[string]deckstate.Zone),
	}

	// Add an empty zone and a non-empty zone
	emptyZone := deckstate.NewZone("empty-zone", deckstate.ZoneTypeTable, 0)
	nonEmptyZone := deckstate.NewZone("non-empty-zone", deckstate.ZoneTypeTable, 0)
	// Add a dummy item to the non-empty zone (using a simple card)
	// Note: This is a mock - in reality we'd need actual card data
	dummyCard := &deckdef.ImageCard{
		ID:   uuid.New(),
		Name: "Dummy Card",
	}
	cardInZone := deckstate.CardInZone{
		Card: dummyCard,
	}
	nonEmptyZone.Items = []deckstate.ZoneItem{cardInZone} // Mock item

	testDeckState.Zones["empty-zone"] = emptyZone
	testDeckState.Zones["non-empty-zone"] = nonEmptyZone

	ctx := context.Background()
	err := mockStateStorage.SaveDeckState(ctx, testDeckState.ID, testDeckState)
	if err != nil {
		t.Fatalf("Failed to save test deck state: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/remove-zone?zone=empty-zone", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleRemoveZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response RemoveZoneResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if !response.Success {
			t.Error("Expected success to be true")
		}

		if response.Meta != nil {
			t.Error("Expected no meta without include=meta query param")
		}
	})

	t.Run("SuccessWithMeta", func(t *testing.T) {
		// First, re-add the zone since it was removed in the previous test
		testDeckState.Zones["empty-zone"] = deckstate.NewZone("empty-zone", deckstate.ZoneTypeTable, 0)
		mockStateStorage.SaveDeckState(ctx, testDeckState.ID, testDeckState)

		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/remove-zone?zone=empty-zone&include=meta", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleRemoveZone(rr, httpReq, "test-state-id")

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var response RemoveZoneResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Errorf("Could not parse response body: %v", err)
		}

		if !response.Success {
			t.Error("Expected success to be true")
		}

		if response.Meta == nil {
			t.Error("Expected meta object when include=meta")
		} else if response.Meta.DurationMS <= 0 {
			t.Errorf("Expected positive duration, got %f", response.Meta.DurationMS)
		}
	})

	t.Run("ZoneNotEmpty", func(t *testing.T) {
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/remove-zone?zone=non-empty-zone", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleRemoveZone(rr, httpReq, "test-state-id")

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

	t.Run("ZoneNotFound", func(t *testing.T) {
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/remove-zone?zone=nonexistent-zone", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleRemoveZone(rr, httpReq, "test-state-id")

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

	t.Run("MissingZoneParameter", func(t *testing.T) {
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/test-state-id/actions/remove-zone", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleRemoveZone(rr, httpReq, "test-state-id")

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

	t.Run("NonExistentDeckState", func(t *testing.T) {
		httpReq, err := http.NewRequest("POST", "/v1/deckstates/nonexistent-id/actions/remove-zone?zone=empty-zone", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.handleRemoveZone(rr, httpReq, "nonexistent-id")

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
}
