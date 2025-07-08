package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

// ZoneRequest represents a request for zone operations (add/remove)
type ZoneRequest struct {
	Name          string              `json:"name"`
	Type          *deckstate.ZoneType `json:"type,omitempty"` // Optional, defaults to ZoneTypeTemporary for add
	DefaultFacing *deckstate.Facing   `json:"default_facing,omitempty"`
	Size          *int                `json:"size,omitempty"`
}

// Validate validates the zone request based on operation type
func (r *ZoneRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	// Validate size hint if provided
	if r.Size != nil && *r.Size < 0 {
		return fmt.Errorf("size must be non-negative")
	}
	return nil
}

func (r *ZoneRequest) Normalize() {
	if r.Type == nil {
		defaultType := deckstate.ZoneTypeTemporary
		r.Type = &defaultType
	}
}

// GetZoneType returns the zone type, defaulting to ZoneTypeTemporary if not specified
func (r *ZoneRequest) GetZoneType() deckstate.ZoneType {
	if r.Type != nil {
		return *r.Type
	}
	return deckstate.ZoneTypeTemporary
}

// ZoneResponse represents a unified response for all zone operations
type ZoneResponse struct {
	Success bool              `json:"success"`
	Zone    *deckstate.Zone   `json:"zone,omitempty"`
	Meta    *ZoneResponseMeta `json:"meta,omitempty"`
}

type ZoneResponseMeta struct {
	Operation  string  `json:"operation"`
	DurationMS float64 `json:"durationMS"`
	// Sort-specific fields
	ZoneLength *int    `json:"zoneLength,omitempty"`
	Sort       *string `json:"sort,omitempty"`
}

func (h *DeckStateHandler) handleAddZone(w http.ResponseWriter, r *http.Request, stateID string) {
	start := time.Now()
	h.logger.Debug("Add zone action requested",
		slog.String("operation", "add_zone"),
		slog.String("deck_state_id", stateID))

	var z ZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&z); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := z.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}
	z.Normalize()

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

	if _, exists := deckState.Zones[z.Name]; exists {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Zone with name '" + z.Name + "' already exists",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	sizeHint := 0
	if z.Size != nil {
		sizeHint = *z.Size
	}

	zone := deckstate.NewZone(z.Name, z.GetZoneType(), sizeHint)

	if z.DefaultFacing != nil {
		zone.DefaultFacing = *z.DefaultFacing
	}

	deckState.Zones[z.Name] = &zone
	if err := h.stateStorage.SaveDeckState(ctx, stateID, deckState); err != nil {
		h.logger.Error("Failed to save deck state after adding zone",
			slog.String("operation", "save_deck_state_after_add_zone"),
			slog.String("deck_state_id", stateID),
			slog.String("zone_name", z.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to save zone",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	duration := time.Since(start)
	zoneResponse := ZoneResponse{
		Success: true,
		Zone:    &zone,
	}

	if r.URL.Query().Get("include") == "meta" {
		zoneResponse.Meta = &ZoneResponseMeta{
			Operation:  "add_zone",
			DurationMS: float64(duration.Microseconds()) / 1000, // Convert to milliseconds
		}
	}

	writeJSONResponse(w, http.StatusCreated, zoneResponse)
}

// handleRemoveZone removes a zone from a deck state
func (h *DeckStateHandler) handleRemoveZone(w http.ResponseWriter, r *http.Request, stateID string) {
	start := time.Now()

	h.logger.Info("Remove zone action requested",
		slog.String("operation", "remove_zone"),
		slog.String("deck_state_id", stateID))

	// Parse request body
	var req ZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := req.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	zoneName := req.Name

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

	delete(deckState.Zones, zoneName)

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
	zoneResponse := ZoneResponse{
		Success: true,
	}

	if r.URL.Query().Get("include") == "meta" {
		zoneResponse.Meta = &ZoneResponseMeta{
			Operation:  "remove_zone",
			DurationMS: float64(duration.Microseconds()) / 1000, // Convert to milliseconds
		}
	}

	writeJSONResponse(w, http.StatusOK, zoneResponse)
}
