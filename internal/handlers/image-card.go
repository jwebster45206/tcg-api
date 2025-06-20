package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// Handler struct with storage dependency
type ImageCardsHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

// NewImageCardsHandler creates a new ImageCardsHandler with the given dependencies
func NewImageCardsHandler(storage storage.Storage, logger *slog.Logger) *ImageCardsHandler {
	return &ImageCardsHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *ImageCardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/image-cards")

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			// GET /image-cards - List all cards
			h.listCards(w, r)
		} else {
			// GET /image-cards/{id} - Get specific card
			cardID := strings.Trim(path, "/")
			h.getCard(w, r, cardID)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			// POST /image-cards - Create new card
			h.createCard(w, r)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodPut:
		if path != "" && path != "/" {
			// PUT /image-cards/{id} - Update card
			cardID := strings.Trim(path, "/")
			h.updateCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for update", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// DELETE /image-cards/{id} - Delete card
			cardID := strings.Trim(path, "/")
			h.deleteCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for deletion", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listCards handles GET /image-cards
func (h *ImageCardsHandler) listCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cards, err := h.storage.ListImageCards(ctx)
	if err != nil {
		h.logger.Error("Failed to list image cards",
			slog.String("operation", "list_image_cards"),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve image cards",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, cards)
}

// getCard handles GET /image-cards/{id}
func (h *ImageCardsHandler) getCard(w http.ResponseWriter, r *http.Request, cardID string) {
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
	card, err := h.storage.GetImageCard(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get image card",
			slog.String("operation", "get_image_card"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "not_found",
			Message: "Card not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, card)
}

// createCard handles POST /image-cards
func (h *ImageCardsHandler) createCard(w http.ResponseWriter, r *http.Request) {
	var card models.ImageCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	createdCard, err := h.storage.CreateImageCard(ctx, card)
	if err != nil {
		h.logger.Error("Failed to create image card",
			slog.String("operation", "create_image_card"),
			slog.String("card_name", card.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create image card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdCard)
}

// updateCard handles PUT /image-cards/{id}
func (h *ImageCardsHandler) updateCard(w http.ResponseWriter, r *http.Request, cardID string) {
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

	var card models.ImageCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	// Set the ID from the URL path
	card.ID = id
	updatedCard, err := h.storage.UpdateImageCard(ctx, card)
	if err != nil {
		h.logger.Error("Failed to update image card",
			slog.String("operation", "update_image_card"),
			slog.String("card_id", cardID),
			slog.String("card_name", card.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update image card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedCard)
}

// deleteCard handles DELETE /image-cards/{id}
func (h *ImageCardsHandler) deleteCard(w http.ResponseWriter, r *http.Request, cardID string) {
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
	if err := h.storage.DeleteImageCard(ctx, id); err != nil {
		h.logger.Error("Failed to delete image card",
			slog.String("operation", "delete_image_card"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete image card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
