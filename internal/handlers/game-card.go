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
type GameCardsHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

// NewGameCardsHandler creates a new GameCardsHandler with the given dependencies
func NewGameCardsHandler(storage storage.Storage, logger *slog.Logger) *GameCardsHandler {
	return &GameCardsHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *GameCardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/game-cards")

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			// GET /game-cards - List all cards
			h.listCards(w, r)
		} else {
			// GET /game-cards/{id} - Get specific card
			cardID := strings.Trim(path, "/")
			h.getCard(w, r, cardID)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			// POST /game-cards - Create new card
			h.createCard(w, r)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodPut:
		if path != "" && path != "/" {
			// PUT /game-cards/{id} - Update card
			cardID := strings.Trim(path, "/")
			h.updateCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for update", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// DELETE /game-cards/{id} - Delete card
			cardID := strings.Trim(path, "/")
			h.deleteCard(w, r, cardID)
		} else {
			http.Error(w, "Card ID required for deletion", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listCards handles GET /game-cards
func (h *GameCardsHandler) listCards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cards, err := h.storage.ListGameCards(ctx, "gamecard")
	if err != nil {
		h.logger.Error("Failed to list cards",
			slog.String("operation", "list_game_cards"),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve cards",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, cards)
}

// getCard handles GET /game-cards/{id}
func (h *GameCardsHandler) getCard(w http.ResponseWriter, r *http.Request, cardID string) {
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
	card, err := h.storage.GetGameCard(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get game card",
			slog.String("operation", "get_game_card"),
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

// createCard handles POST /game-cards
func (h *GameCardsHandler) createCard(w http.ResponseWriter, r *http.Request) {
	var card models.GameCard

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		response := ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	createdCard, err := h.storage.CreateGameCard(ctx, card)
	if err != nil {
		h.logger.Error("Failed to create game card",
			slog.String("operation", "create_game_card"),
			slog.String("card_name", card.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusCreated, createdCard)
}

// updateCard handles PUT /game-cards/{id}
func (h *GameCardsHandler) updateCard(w http.ResponseWriter, r *http.Request, cardID string) {
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

	var card models.GameCard
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
	updatedCard, err := h.storage.UpdateGameCard(ctx, card)
	if err != nil {
		h.logger.Error("Failed to update game card",
			slog.String("operation", "update_game_card"),
			slog.String("card_id", cardID),
			slog.String("card_name", card.Name),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeJSONResponse(w, http.StatusOK, updatedCard)
}

// deleteCard handles DELETE /game-cards/{id}
func (h *GameCardsHandler) deleteCard(w http.ResponseWriter, r *http.Request, cardID string) {
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
	if err := h.storage.DeleteGameCard(ctx, id); err != nil {
		h.logger.Error("Failed to delete game card",
			slog.String("operation", "delete_game_card"),
			slog.String("card_id", cardID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete card",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
