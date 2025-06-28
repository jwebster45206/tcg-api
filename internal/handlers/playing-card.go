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
type PlayingCardsHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

// NewPlayingCardsHandler creates a new PlayingCardsHandler with the given dependencies
func NewPlayingCardsHandler(storage storage.Storage, logger *slog.Logger) *PlayingCardsHandler {
	return &PlayingCardsHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *PlayingCardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/playing-cards")

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			// GET /playing-cards - List all cards
			h.listCards(w, r)
		} else {
			// GET /playing-cards/{id} - Get specific card
			cardID := strings.Trim(path, "/")
			h.getCard(w, r, cardID)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			// POST /playing-cards - Create new card
			h.createCard(w, r)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodPatch:
		if path != "" && path != "/" {
			// PATCH /playing-cards/{id} - Update card
			cardID := strings.Trim(path, "/")
			h.updateCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for update", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// DELETE /playing-cards/{id} - Delete card
			cardID := strings.Trim(path, "/")
			h.deleteCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for deletion", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listCards handles GET /playing-cards
func (h *PlayingCardsHandler) listCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters, err := ParseFilters(r, models.PlayingCardQueryConfig)
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

	sorts, err := ParseSorts(r, models.PlayingCardQueryConfig)
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
	pageSize := limit
	pageNum := (offset / limit) + 1

	cards, err := h.storage.ListPlayingCards(ctx, filters, sorts, pageSize, pageNum)
	if err != nil {
		h.logger.Error("Failed to list playing cards",
			slog.String("operation", "list_playing_cards"),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to retrieve playing cards",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, cards)
}

// getCard handles GET /playing-cards/{id}
func (h *PlayingCardsHandler) getCard(w http.ResponseWriter, r *http.Request, cardID string) {
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
	card, err := h.storage.GetPlayingCard(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get playing card",
			slog.String("operation", "get_playing_card"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Playing card not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, card)
}

// createCard handles POST /playing-cards
func (h *PlayingCardsHandler) createCard(w http.ResponseWriter, r *http.Request) {
	var card models.PlayingCard

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate the playing card
	if err := card.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrValidation,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	createdCard, err := h.storage.CreatePlayingCard(ctx, card)
	if err != nil {
		h.logger.Error("Failed to create playing card",
			slog.String("operation", "create_playing_card"),
			slog.String("suit", card.Suit),
			slog.Int("ranking", card.Ranking),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to create playing card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdCard)
}

// updateCard handles PUT/PATCH /playing-cards/{id}
func (h *PlayingCardsHandler) updateCard(w http.ResponseWriter, r *http.Request, cardID string) {
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

	// Get the existing card first
	existingCard, err := h.storage.GetPlayingCard(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get existing playing card for update",
			slog.String("operation", "get_playing_card_for_update"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Playing card not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	// For PATCH, decode partial updates onto existing card
	// For PUT, decode complete replacement
	var updateCard models.PlayingCard
	if r.Method == http.MethodPatch {
		// Start with existing card data
		updateCard = *existingCard
	}

	if err := json.NewDecoder(r.Body).Decode(&updateCard); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Set the ID from the URL path (ensure it doesn't get overridden)
	updateCard.ID = id

	// Validate the updated card
	if err := updateCard.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrValidation,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	updatedCard, err := h.storage.UpdatePlayingCard(ctx, updateCard)
	if err != nil {
		h.logger.Error("Failed to update playing card",
			slog.String("operation", "update_playing_card"),
			slog.String("card_id", cardID),
			slog.String("suit", updateCard.Suit),
			slog.Int("ranking", updateCard.Ranking),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to update playing card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedCard)
}

// deleteCard handles DELETE /playing-cards/{id}
func (h *PlayingCardsHandler) deleteCard(w http.ResponseWriter, r *http.Request, cardID string) {
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
	if err := h.storage.DeletePlayingCard(ctx, id); err != nil {
		h.logger.Error("Failed to delete playing card",
			slog.String("operation", "delete_playing_card"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to delete playing card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
