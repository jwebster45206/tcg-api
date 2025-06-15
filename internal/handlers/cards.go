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
			h.listCards(w, r)
		} else {
			// GET /cards/{id} - Get specific card
			cardID := strings.Trim(path, "/")
			h.getCard(w, r, cardID)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			// POST /cards - Create new card
			h.createCard(w, r)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodPut:
		if path != "" && path != "/" {
			// PUT /cards/{id} - Update card
			cardID := strings.Trim(path, "/")
			h.updateCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for update", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// DELETE /cards/{id} - Delete card
			cardID := strings.Trim(path, "/")
			h.deleteCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for deletion", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listCards handles GET /cards
func (h *CardsHandler) listCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cards, err := h.storage.ListCards(ctx)
	if err != nil {
		h.logger.Printf("Failed to list cards: %v", err)
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve cards",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, cards)
}

// getCard handles GET /cards/{id}
func (h *CardsHandler) getCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	id, err := uuid.Parse(cardID)
	if err != nil {
		response := ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	card, err := h.storage.GetCard(ctx, id)
	if err != nil {
		h.logger.Printf("Failed to get card %s: %v", cardID, err)
		response := ErrorResponse{
			Error:   "not_found",
			Message: "Card not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, card)
}

// createCard handles POST /cards
func (h *CardsHandler) createCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	if err := h.storage.CreateCard(ctx, &card); err != nil {
		h.logger.Printf("Failed to create card: %v", err)
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusCreated, card)
}

// updateCard handles PUT /cards/{id}
func (h *CardsHandler) updateCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	id, err := uuid.Parse(cardID)
	if err != nil {
		response := ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var card models.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	if err := h.storage.UpdateCard(ctx, id, &card); err != nil {
		h.logger.Printf("Failed to update card %s: %v", cardID, err)
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, card)
}

// deleteCard handles DELETE /cards/{id}
func (h *CardsHandler) deleteCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	id, err := uuid.Parse(cardID)
	if err != nil {
		response := ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	if err := h.storage.DeleteCard(ctx, id); err != nil {
		h.logger.Printf("Failed to delete card %s: %v", cardID, err)
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
