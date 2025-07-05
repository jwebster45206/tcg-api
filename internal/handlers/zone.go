package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

// AddZoneRequest represents the request to add a new zone
type AddZoneRequest struct {
	Name          string             `json:"name"`
	Type          deckstate.ZoneType `json:"type"`
	DefaultFacing *deckstate.Facing  `json:"default_facing,omitempty"`
	Size          *int               `json:"size,omitempty"` // Optional size hint for initial capacity
}

// AddZoneResponse represents the response from adding a zone
type AddZoneResponse struct {
	Zone *deckstate.Zone `json:"zone"`
	Meta *AddZoneMeta    `json:"meta,omitempty"`
}

// AddZoneMeta contains metadata about the add zone operation
type AddZoneMeta struct {
	DurationMS float64 `json:"durationMS"`
}

// RemoveZoneResponse represents the response from removing a zone
type RemoveZoneResponse struct {
	Success bool            `json:"success"`
	Meta    *RemoveZoneMeta `json:"meta,omitempty"`
}

// RemoveZoneMeta contains metadata about the remove zone operation
type RemoveZoneMeta struct {
	DurationMS float64 `json:"durationMS"`
}

// handleAddZone adds a new zone to a deck state
func (h *DeckStateHandler) handleAddZone(w http.ResponseWriter, r *http.Request, stateID string) {
	start := time.Now()

	h.logger.Info("Add zone action requested",
		slog.String("operation", "add_zone"),
		slog.String("deck_state_id", stateID))

	// Parse request body
	var req AddZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate required fields
	if req.Name == "" {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "name is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if req.Type == "" {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "type is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	deckState, err := h.stateStorage.GetDeckState(ctx, stateID)
	if err != nil {
		h.logger.Error("Failed to get deck state for add zone",
			slog.String("operation", "get_deck_state_for_add_zone"),
			slog.String("deck_state_id", stateID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to retrieve deck state",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if deckState == nil {
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Deck state not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	// Check if zone name already exists
	if _, exists := deckState.Zones[req.Name]; exists {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Zone with name '" + req.Name + "' already exists",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Determine size hint (default to 0 for unlimited)
	sizeHint := 0
	if req.Size != nil {
		sizeHint = *req.Size
		if sizeHint < 0 {
			response := ErrorResponse{
				Error:   errStrBadRequest,
				Message: "size_hint must be non-negative",
			}
			writeJSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}

	// Create new zone
	zone := deckstate.NewZone(req.Name, req.Type, sizeHint)

	// Override default facing if provided
	if req.DefaultFacing != nil {
		zone.DefaultFacing = *req.DefaultFacing
	}

	// Add zone to deck state
	deckState.Zones[req.Name] = zone

	// Save the updated deck state
	if err := h.stateStorage.SaveDeckState(ctx, stateID, deckState); err != nil {
		h.logger.Error("Failed to save deck state after adding zone",
			slog.String("operation", "save_deck_state_after_add_zone"),
			slog.String("deck_state_id", stateID),
			slog.String("zone_name", req.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to save zone",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	duration := time.Since(start)

	h.logger.Info("Successfully added zone",
		slog.String("operation", "add_zone"),
		slog.String("deck_state_id", stateID),
		slog.String("zone_name", req.Name),
		slog.String("zone_type", string(req.Type)),
		slog.Duration("duration", duration))

	// Prepare response
	addZoneResponse := AddZoneResponse{
		Zone: &zone,
	}

	// Include meta if requested via query parameter
	if r.URL.Query().Get("include") == "meta" {
		addZoneResponse.Meta = &AddZoneMeta{
			DurationMS: float64(duration.Microseconds()) / 1000, // Convert to milliseconds
		}
	}

	writeJSONResponse(w, http.StatusCreated, addZoneResponse)
}

// handleRemoveZone removes a zone from a deck state
func (h *DeckStateHandler) handleRemoveZone(w http.ResponseWriter, r *http.Request, stateID string) {
	start := time.Now()

	h.logger.Info("Remove zone action requested",
		slog.String("operation", "remove_zone"),
		slog.String("deck_state_id", stateID))

	// Get zone name from query parameter
	zoneName := r.URL.Query().Get("zone")
	if zoneName == "" {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "zone parameter is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	deckState, err := h.stateStorage.GetDeckState(ctx, stateID)
	if err != nil {
		h.logger.Error("Failed to get deck state for remove zone",
			slog.String("operation", "get_deck_state_for_remove_zone"),
			slog.String("deck_state_id", stateID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to retrieve deck state",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	if deckState == nil {
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Deck state not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	// Check if zone exists
	zone, exists := deckState.Zones[zoneName]
	if !exists {
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Zone '" + zoneName + "' not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	// Check if zone is empty
	if len(zone.Items) > 0 {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Cannot remove zone '" + zoneName + "' because it contains " + fmt.Sprintf("%d", len(zone.Items)) + " item(s)",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Remove zone from deck state
	delete(deckState.Zones, zoneName)

	// Save the updated deck state
	if err := h.stateStorage.SaveDeckState(ctx, stateID, deckState); err != nil {
		h.logger.Error("Failed to save deck state after removing zone",
			slog.String("operation", "save_deck_state_after_remove_zone"),
			slog.String("deck_state_id", stateID),
			slog.String("zone_name", zoneName),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to remove zone",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	duration := time.Since(start)

	h.logger.Info("Successfully removed zone",
		slog.String("operation", "remove_zone"),
		slog.String("deck_state_id", stateID),
		slog.String("zone_name", zoneName),
		slog.Duration("duration", duration))

	// Prepare response
	removeZoneResponse := RemoveZoneResponse{
		Success: true,
	}

	// Include meta if requested via query parameter
	if r.URL.Query().Get("include") == "meta" {
		removeZoneResponse.Meta = &RemoveZoneMeta{
			DurationMS: float64(duration.Microseconds()) / 1000, // Convert to milliseconds
		}
	}

	writeJSONResponse(w, http.StatusOK, removeZoneResponse)
}
