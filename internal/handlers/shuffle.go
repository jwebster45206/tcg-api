package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jwebster45206/tcg-api/internal/deckstate"
	"github.com/jwebster45206/tcg-api/internal/shuffle"
)

const (
	SortTypeShuffle = "shuffle"
)

// SortZoneRequest represents the request to sort a zone
type SortZoneRequest struct {
	Zone string `json:"zone"`
	Sort string `json:"sort"`
}

// SortZoneResponse represents the response from a sort zone operation
type SortZoneResponse struct {
	Zone *deckstate.Zone `json:"zone"`
	Meta *SortZoneMeta   `json:"meta,omitempty"`
}

// SortZoneMeta contains metadata about the sort operation
type SortZoneMeta struct {
	ZoneLength int    `json:"zoneLength"`
	Sort       string `json:"sort"`
	DurationMS int64  `json:"durationMS"`
}

// shuffleZone performs a shuffle operation on a zone, returning
// a measurement of the time taken to perform the shuffle.
func shuffleZone(zone *deckstate.Zone) (time.Duration, error) {
	start := time.Now()
	err := shuffle.FisherYatesShuffle(zone.Items)
	if err != nil {
		return 0, err
	}
	duration := time.Since(start)
	return duration, nil
}

// handleSortZone sorts cards within a zone
func (h *DeckStateHandler) handleSortZone(w http.ResponseWriter, r *http.Request, stateID string) {
	h.logger.Info("Sort zone action requested",
		slog.String("operation", "sort_zone"),
		slog.String("deck_state_id", stateID))

	// Parse request body
	var req SortZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate required fields
	if req.Zone == "" {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "zone is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if req.Sort == "" {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "sort is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if req.Sort != SortTypeShuffle {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "only 'shuffle' sort type is supported",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	deckState, err := h.stateStorage.GetDeckState(ctx, stateID)
	if err != nil {
		h.logger.Error("Failed to get deck state for sort",
			slog.String("operation", "get_deck_state_for_sort"),
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

	zone, exists := deckState.Zones[req.Zone]
	if !exists {
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Zone not found: " + req.Zone,
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	// Perform the shuffle operation
	duration, err := shuffleZone(&zone)
	if err != nil {
		h.logger.Error("Failed to shuffle zone",
			slog.String("operation", "shuffle_zone"),
			slog.String("deck_state_id", stateID),
			slog.String("zone", req.Zone),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to shuffle zone",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	deckState.Zones[req.Zone] = zone
	if err := h.stateStorage.SaveDeckState(ctx, stateID, deckState); err != nil {
		h.logger.Error("Failed to save deck state after shuffle",
			slog.String("operation", "save_deck_state_after_shuffle"),
			slog.String("deck_state_id", stateID),
			slog.String("zone", req.Zone),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to save shuffled zone",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	query := r.URL.Query()
	includes := strings.Split(query.Get("include"), ",")

	responseZone := zone
	includeItems := false
	includeMeta := false
	for _, include := range includes {
		switch strings.TrimSpace(include) {
		case "items":
			includeItems = true
		case "meta":
			includeMeta = true
		}
	}

	// If items are not explicitly included, remove them
	if !includeItems {
		responseZone.Items = nil
	}

	sortResponse := SortZoneResponse{
		Zone: &responseZone,
	}

	if includeMeta {
		sortResponse.Meta = &SortZoneMeta{
			ZoneLength: len(zone.Items),
			Sort:       req.Sort,
			DurationMS: duration.Milliseconds(),
		}
	}

	writeJSONResponse(w, http.StatusOK, sortResponse)
}
