package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jwebster45206/tcg-api/internal/shuffle"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
	"github.com/jwebster45206/tcg-api/pkg/deckstate"
)

const (
	SortTypeShuffle    = "shuffle"
	SortTypeDefinition = "definition"
)

// SortZoneRequest represents the request to sort a zone
type SortZoneRequest struct {
	Zone string `json:"zone"`
	Sort string `json:"sort"`
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

// sortZoneByDefinition performs a definition sort operation on a zone, returning
// a measurement of the time taken to perform the sort.
func sortZoneByDefinition(zone *deckstate.Zone, deck *deckdef.Deck) (time.Duration, error) {
	start := time.Now()
	err := shuffle.DefinitionSort(zone.Items, deck)
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

	if req.Sort != SortTypeShuffle && req.Sort != SortTypeDefinition {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "only 'shuffle' and 'definition' sort types are supported",
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

	// Perform the sort operation
	var duration time.Duration

	switch req.Sort {
	case SortTypeShuffle:
		// Fisher-Yates randomization of the zone
		duration, err = shuffleZone(zone)
	case SortTypeDefinition:
		// Reset to deck definition
		duration, err = sortZoneByDefinition(zone, &deckState.Deck)
	}

	if err != nil {
		h.logger.Error("Failed to sort zone",
			slog.String("operation", "sort_zone"),
			slog.String("deck_state_id", stateID),
			slog.String("zone", req.Zone),
			slog.String("sort_type", req.Sort),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to sort zone",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// No need to reassign zone since the pointer is stored in the map
	if err := h.stateStorage.SaveDeckState(ctx, stateID, deckState); err != nil {
		h.logger.Error("Failed to save deck state after sort",
			slog.String("operation", "save_deck_state_after_sort"),
			slog.String("deck_state_id", stateID),
			slog.String("zone", req.Zone),
			slog.String("sort_type", req.Sort),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to save sorted zone",
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

	sortResponse := ZoneResponse{
		Success: true,
		Zone:    responseZone,
	}

	if includeMeta {
		zoneLength := len(zone.Items)
		sortResponse.Meta = &ZoneResponseMeta{
			Operation:  "sort_zone",
			DurationMS: float64(duration.Microseconds()) / 1000, // Convert to milliseconds
			ZoneLength: &zoneLength,
			Sort:       &req.Sort,
		}
	}

	writeJSONResponse(w, http.StatusOK, sortResponse)
}
