package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

type DrawRequest struct {
	FromZone string `json:"from_zone"`
	ToZone   string `json:"to_zone"`
	Count    int    `json:"count"`
}

func (r *DrawRequest) Validate() error {
	if r.FromZone == "" {
		return fmt.Errorf("from_zone is required")
	}
	if r.ToZone == "" {
		return fmt.Errorf("to_zone is required")
	}
	if r.Count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}
	return nil
}

type DrawResponse struct {
	Success bool                       `json:"success"`
	Zones   map[string]*deckstate.Zone `json:"zones,omitempty"`
	Meta    *DrawResponseMeta          `json:"meta,omitempty"`
}

type DrawResponseMeta struct {
	Operation  string  `json:"operation"`
	DurationMS float64 `json:"durationMS"`
	FromZone   string  `json:"from_zone"`
	ToZone     string  `json:"to_zone"`
	CardsDrawn int     `json:"cards_drawn"`
}

// handleDrawCards moves cards from the end of one zone to the end of another
func (h *DeckStateHandler) handleDrawCards(w http.ResponseWriter, r *http.Request, stateID string) {
	start := time.Now()
	h.logger.Info("Draw cards action requested",
		slog.String("operation", "draw_cards"),
		slog.String("deck_state_id", stateID))

	var req DrawRequest
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

	ctx := r.Context()
	deckState, err := h.stateStorage.GetDeckState(ctx, stateID)
	if err != nil {
		h.logger.Error("Failed to get deck state for draw",
			slog.String("operation", "get_deck_state_for_draw"),
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

	// Note: Shallow copys of zones
	fromZone, exists := deckState.Zones[req.FromZone]
	if !exists {
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Source zone not found: " + req.FromZone,
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	toZone, exists := deckState.Zones[req.ToZone]
	if !exists {
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Destination zone not found: " + req.ToZone,
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	if len(fromZone.Items) == 0 {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Source zone is empty",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}
	if len(fromZone.Items) < req.Count {
		req.Count = len(fromZone.Items)
	}

	// Pop cards from the end of the source zone
	sourceLen := len(fromZone.Items)
	cardsToDraw := fromZone.Items[sourceLen-req.Count:]
	fromZone.Items = fromZone.Items[:sourceLen-req.Count]

	// Push cards to the end of the destination zone
	toZone.Items = append(toZone.Items, cardsToDraw...)

	// Save the updated deck state
	// No need to reassign zones singe pointers are stored in the map
	if err := h.stateStorage.SaveDeckState(ctx, stateID, deckState); err != nil {
		h.logger.Error("Failed to save deck state after draw",
			slog.String("operation", "save_deck_state_after_draw"),
			slog.String("deck_state_id", stateID),
			slog.String("from_zone", req.FromZone),
			slog.String("to_zone", req.ToZone),
			slog.Int("count", req.Count),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to save deck state after drawing cards",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	duration := time.Since(start)

	query := r.URL.Query()
	includes := strings.Split(query.Get("include"), ",")

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

	// Prepare response zones
	responseZones := make(map[string]*deckstate.Zone)

	// Create copies of the zones for the response
	// Note: Since we're using pointers, we need to copy the struct values
	fromZoneResponse := *fromZone
	toZoneResponse := *toZone

	// If items are not explicitly included, remove them
	if !includeItems {
		fromZoneResponse.Items = nil
		toZoneResponse.Items = nil
	}

	responseZones[req.FromZone] = &fromZoneResponse
	responseZones[req.ToZone] = &toZoneResponse

	drawResponse := DrawResponse{
		Success: true,
		Zones:   responseZones,
	}

	if includeMeta {
		drawResponse.Meta = &DrawResponseMeta{
			Operation:  "draw_cards",
			DurationMS: float64(duration.Microseconds()) / 1000, // Convert to milliseconds
			FromZone:   req.FromZone,
			ToZone:     req.ToZone,
			CardsDrawn: req.Count,
		}
	}

	writeJSONResponse(w, http.StatusOK, drawResponse)
}
