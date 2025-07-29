package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/storage"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
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
	path := strings.TrimPrefix(r.URL.Path, "/v1/image-cards")

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

	case http.MethodPatch:
		if path != "" && path != "/" {
			// PATCH /image-cards/{id} - Update card
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

	filters, err := ParseFilters(r, deckdef.ImageCardQueryConfig)
	if err != nil {
		h.logger.Error("Failed to parse filters",
			slog.String("operation", "parse_filters"),
			slog.Any("error", err))

		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Parse sorts from query parameters
	sorts, err := ParseSorts(r, deckdef.ImageCardQueryConfig)
	if err != nil {
		h.logger.Error("Failed to parse sorts",
			slog.String("operation", "parse_sorts"),
			slog.Any("error", err))

		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Parse pagination from query parameters
	offset, limit, err := ParsePagination(r)
	if err != nil {
		h.logger.Error("Failed to parse pagination",
			slog.String("operation", "parse_pagination"),
			slog.Any("error", err))

		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Convert offset/limit to page-based parameters for storage layer
	pageSize := limit
	pageNum := (offset / limit) + 1

	cards, err := h.storage.ListImageCards(ctx, filters, sorts, pageSize, pageNum)
	if err != nil {
		h.logger.Error("Failed to list image cards",
			slog.String("operation", "list_image_cards"),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
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
			Error:   errStrBadRequest,
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
	var card deckdef.ImageCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate the card using the model's validation method
	if err := card.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
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
			Error:   errStrBadRequest,
			Message: "Failed to create image card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdCard)
}

// updateCard handles PATCH /image-cards/{id}
func (h *ImageCardsHandler) updateCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	id, err := uuid.Parse(cardID)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid card ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()

	// Fetch the current card
	currentCard, err := h.storage.GetImageCard(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get current image card for update",
			slog.String("operation", "get_image_card_for_update"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Card not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	// Parse the update request into a temporary struct
	var updateData struct {
		Name          *string `json:"name"`
		Description   *string `json:"description"`
		FrontImageURL *string `json:"front_image_url"`
		BackImageURL  *string `json:"back_image_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Apply updates only for fields that were provided
	updatedCard := *currentCard // Copy the current card

	if updateData.Name != nil {
		updatedCard.Name = *updateData.Name
	}
	if updateData.Description != nil {
		updatedCard.Description = *updateData.Description
	}
	if updateData.FrontImageURL != nil {
		updatedCard.FrontImageURL = *updateData.FrontImageURL
	}
	if updateData.BackImageURL != nil {
		updatedCard.BackImageURL = *updateData.BackImageURL
	}

	// Validate the updated card
	if err := updatedCard.Validate(); err != nil {
		response := ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Save the updated card
	savedCard, err := h.storage.UpdateImageCard(ctx, updatedCard)
	if err != nil {
		h.logger.Error("Failed to update image card",
			slog.String("operation", "update_image_card"),
			slog.String("card_id", cardID),
			slog.String("card_name", updatedCard.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to update image card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, savedCard)
}

// deleteCard handles DELETE /image-cards/{id}
func (h *ImageCardsHandler) deleteCard(w http.ResponseWriter, r *http.Request, cardID string) {
	// Validate UUID format
	id, err := uuid.Parse(cardID)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
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
			Error:   errStrInternal,
			Message: "Failed to delete image card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
