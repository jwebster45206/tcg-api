package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// Handler struct with storage dependency
type CardsHandler struct {
	storage storage.Storage
	logger  *log.Logger
}

// NewCardsHandler creates a new CardsHandler with the given dependencies
func NewCardsHandler(storage storage.Storage, logger *log.Logger) *CardsHandler {
	return &CardsHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *CardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/cards")

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			// GET /cards - List all cards
			listCards(w, r)
		} else {
			// GET /cards/{id} - Get specific card
			cardID := strings.Trim(path, "/")
			getCard(w, r, cardID)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			// POST /cards - Create new card
			createCard(w, r)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodPut:
		if path != "" && path != "/" {
			// PUT /cards/{id} - Update card
			cardID := strings.Trim(path, "/")
			updateCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for update", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// DELETE /cards/{id} - Delete card
			cardID := strings.Trim(path, "/")
			deleteCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for deletion", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listCards handles GET /cards
func listCards(w http.ResponseWriter, r *http.Request) {
	response := ErrorResponse{
		Error:   "not_implemented",
		Message: "List cards endpoint not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}

// getCard handles GET /cards/{id}
func getCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	if _, err := uuid.Parse(cardID); err != nil {
		response := ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	response := ErrorResponse{
		Error:   "not_implemented",
		Message: "Get card endpoint not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}

// createCard handles POST /cards
func createCard(w http.ResponseWriter, r *http.Request) {
	var c models.Card

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	response := ErrorResponse{
		Error:   "not_implemented",
		Message: "Create card endpoint not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}

// updateCard handles PUT /cards/{id}
func updateCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	if _, err := uuid.Parse(cardID); err != nil {
		response := ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var c models.Card

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	response := ErrorResponse{
		Error:   "not_implemented",
		Message: "Update card endpoint not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}

// deleteCard handles DELETE /cards/{id}
func deleteCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	if _, err := uuid.Parse(cardID); err != nil {
		response := ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	response := ErrorResponse{
		Error:   "not_implemented",
		Message: "Delete card endpoint not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}
