package handlers

import (
	"log/slog"
	"net/http"
)

// handleSortZone sorts cards within a zone
func (h *DeckStateHandler) handleSortZone(w http.ResponseWriter, r *http.Request, stateID string) {
	h.logger.Info("Sort zone action requested",
		slog.String("operation", "sort_zone"),
		slog.String("deck_state_id", stateID))

	response := ErrorResponse{
		Error:   errStrNotImplemented,
		Message: "Sort zone action not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}
