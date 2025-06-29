package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// DecksHandler handles HTTP requests for deck operations
type DecksHandler struct {
	storage storage.Storage
	logger  *slog.Logger
}

// NewDecksHandler creates a new DecksHandler with the given dependencies
func NewDecksHandler(storage storage.Storage, logger *slog.Logger) *DecksHandler {
	return &DecksHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *DecksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/decks")

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			// GET /decks - List all decks
			h.listDecks(w, r)
		} else {
			// GET /decks/{id} - Get specific deck
			deckID := strings.Trim(path, "/")
			h.getDeck(w, r, deckID)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			// POST /decks - Create new deck
			h.createDeck(w, r)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodPatch:
		if path != "" && path != "/" {
			// PATCH /decks/{id} - Update specific deck
			deckID := strings.Trim(path, "/")
			h.updateDeck(w, r, deckID)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// DELETE /decks/{id} - Delete specific deck
			deckID := strings.Trim(path, "/")
			h.deleteDeck(w, r, deckID)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listDecks handles GET /decks with filtering, sorting, and pagination
func (h *DecksHandler) listDecks(w http.ResponseWriter, r *http.Request) {
	// Parse filters
	filters, err := ParseFilters(r, models.DeckQueryConfig)
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

	sorts, err := ParseSorts(r, models.DeckQueryConfig)
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

	pageSize, pageNum, err := ParsePagination(r)
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
	decks, err := h.storage.ListDecks(r.Context(), filters, sorts, pageSize, pageNum)
	if err != nil {
		h.logger.Error("Failed to list decks",
			slog.String("operation", "list_decks"),
			slog.Any("error", err))

		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to retrieve decks",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Check if cards should be included
	includeCards := shouldIncludeCards(r.URL.Query().Get("include"))

	// Only allow including cards if filtering by ID (to avoid expensive operations on large lists)
	if includeCards {
		hasIDFilter := false
		for _, filter := range filters {
			if filter.Column == "id" {
				hasIDFilter = true
				break
			}
		}
		if !hasIDFilter {
			response := ErrorResponse{
				Error:   errStrBadRequest,
				Message: "Including cards is only allowed when filtering by 'id'",
			}
			writeJSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}

	// If cards are requested, fetch and populate them for each deck
	if includeCards {
		for _, deck := range decks {
			cards, err := h.storage.ListDeckCards(r.Context(), deck.ID)
			if err != nil {
				h.logger.Error("Failed to get deck cards",
					slog.String("operation", "get_deck_cards"),
					slog.String("deck_id", deck.ID.String()),
					slog.Any("error", err))
				// Continue with other decks instead of failing completely
				continue
			}

			// Build CardCollection
			totalCount := 0
			for _, cardWithQuantity := range cards {
				totalCount += cardWithQuantity.Quantity
			}

			cardCollection := &models.CardCollection{
				TotalCount:  totalCount,
				UniqueCount: len(cards),
				Items:       make([]models.CardWithQuantity, len(cards)),
			}

			// Convert from pointers to values for the response
			for i, cardPtr := range cards {
				cardCollection.Items[i] = *cardPtr
			}

			deck.Cards = cardCollection
		}
	}

	// Set content type and return JSON
	writeJSONResponse(w, http.StatusOK, decks)
}

// getDeck handles GET /decks/{id}
func (h *DecksHandler) getDeck(w http.ResponseWriter, r *http.Request, deckIDStr string) {
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid deck ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	deck, err := h.storage.GetDeck(r.Context(), deckID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			response := ErrorResponse{
				Error:   errStrNotFound,
				Message: "Deck not found",
			}
			writeJSONResponse(w, http.StatusNotFound, response)
		} else {
			h.logger.Error("Failed to get deck",
				slog.String("operation", "get_deck"),
				slog.String("deck_id", deckIDStr),
				slog.Any("error", err))

			response := ErrorResponse{
				Error:   errStrInternal,
				Message: "Failed to retrieve deck",
			}
			writeJSONResponse(w, http.StatusInternalServerError, response)
		}
		return
	}

	// If cards are requested, fetch and populate them
	includeCards := shouldIncludeCards(r.URL.Query().Get("include"))
	if includeCards {
		cards, err := h.storage.ListDeckCards(r.Context(), deckID)
		if err != nil {
			h.logger.Error("Failed to get deck cards",
				slog.String("operation", "get_deck_cards"),
				slog.String("deck_id", deckIDStr),
				slog.Any("error", err))

			response := ErrorResponse{
				Error:   errStrInternal,
				Message: "Failed to retrieve deck cards",
			}
			writeJSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		totalCount := 0
		for _, cardWithQuantity := range cards {
			totalCount += cardWithQuantity.Quantity
		}

		cardCollection := &models.CardCollection{
			TotalCount:  totalCount,
			UniqueCount: len(cards),
			Items:       make([]models.CardWithQuantity, len(cards)),
		}

		// Convert from pointers to values for the response
		for i, cardPtr := range cards {
			cardCollection.Items[i] = *cardPtr
		}
		deck.Cards = cardCollection
	}

	writeJSONResponse(w, http.StatusOK, deck)
}

// createDeck handles POST /decks
func (h *DecksHandler) createDeck(w http.ResponseWriter, r *http.Request) {
	var deckInput models.DeckInput
	if err := json.NewDecoder(r.Body).Decode(&deckInput); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate the input using the model's validation method
	if err := deckInput.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Convert input to deck model for basic field updates
	deck := models.Deck{
		Name:           deckInput.Name,
		Type:           deckInput.DeckType,
		SleeveImageURL: deckInput.SleeveImageURL,
	}

	createdDeck, err := h.storage.CreateDeck(r.Context(), deck)
	if err != nil {
		h.logger.Error("Failed to create deck",
			slog.String("operation", "create_deck"),
			slog.String("deck_name", deckInput.Name),
			slog.Any("error", err))

		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to create deck",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// If cards are provided, set them for the new deck
	if deckInput.Cards != nil && len(deckInput.Cards.Items) > 0 {
		err = h.storage.SetDeckCards(r.Context(), createdDeck.ID, deckInput.Cards.Items)
		if err != nil {
			h.logger.Error("Failed to set deck cards",
				slog.String("operation", "set_deck_cards"),
				slog.String("deck_id", createdDeck.ID.String()),
				slog.Any("error", err))
			// Note: We could consider rolling back the deck creation here, but for now we'll just log the error
			// The deck exists but without cards - client can retry setting cards later
		}
	}

	// Check if response should include cards
	includeCards := shouldIncludeCards(r.URL.Query().Get("include"))
	if includeCards {
		cards, err := h.storage.ListDeckCards(r.Context(), createdDeck.ID)
		if err != nil {
			h.logger.Error("Failed to get deck cards for response",
				slog.String("operation", "get_deck_cards"),
				slog.String("deck_id", createdDeck.ID.String()),
				slog.Any("error", err))
			// Continue without cards in response rather than failing the whole request
		} else {
			totalCount := 0
			for _, cardWithQuantity := range cards {
				totalCount += cardWithQuantity.Quantity
			}

			cardCollection := &models.CardCollection{
				TotalCount:  totalCount,
				UniqueCount: len(cards),
				Items:       make([]models.CardWithQuantity, len(cards)),
			}

			// Convert from pointers to values for the response
			for i, cardPtr := range cards {
				cardCollection.Items[i] = *cardPtr
			}
			createdDeck.Cards = cardCollection
		}
	}

	writeJSONResponse(w, http.StatusCreated, createdDeck)
}

// updateDeck handles PATCH /decks/{id}
func (h *DecksHandler) updateDeck(w http.ResponseWriter, r *http.Request, deckIDStr string) {
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid deck ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Parse request body
	var deckInput models.DeckInput
	if err := json.NewDecoder(r.Body).Decode(&deckInput); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate the input using the model's validation method
	if err := deckInput.Validate(); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: err.Error(),
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Convert input to deck model for basic field updates
	deck := models.Deck{
		ID:             deckID,
		Name:           deckInput.Name,
		Type:           deckInput.DeckType,
		SleeveImageURL: deckInput.SleeveImageURL,
	}

	// Update the deck's basic information
	updatedDeck, err := h.storage.UpdateDeck(r.Context(), deck)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			response := ErrorResponse{
				Error:   errStrNotFound,
				Message: "Deck not found",
			}
			writeJSONResponse(w, http.StatusNotFound, response)
		} else {
			h.logger.Error("Failed to update deck",
				slog.String("operation", "update_deck"),
				slog.String("deck_id", deckIDStr),
				slog.Any("error", err))

			response := ErrorResponse{
				Error:   errStrInternal,
				Message: "Failed to update deck",
			}
			writeJSONResponse(w, http.StatusInternalServerError, response)
		}
		return
	}

	// Handle card updates if provided
	if deckInput.Cards != nil {
		// Get current deck cards to compare
		currentCards, err := h.storage.ListDeckCards(r.Context(), deckID)
		if err != nil {
			h.logger.Error("Failed to get current deck cards",
				slog.String("operation", "get_current_deck_cards"),
				slog.String("deck_id", deckIDStr),
				slog.Any("error", err))

			response := ErrorResponse{
				Error:   errStrInternal,
				Message: "Failed to retrieve current deck cards",
			}
			writeJSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		// Use our efficient comparison function to check if cards changed
		if deckCardsChangedFromInput(currentCards, deckInput.Cards) {
			err = h.storage.SetDeckCards(r.Context(), deckID, deckInput.Cards.Items)
			if err != nil {
				h.logger.Error("Failed to update deck cards",
					slog.String("operation", "update_deck_cards"),
					slog.String("deck_id", deckIDStr),
					slog.Any("error", err))

				response := ErrorResponse{
					Error:   errStrInternal,
					Message: "Failed to update deck cards",
				}
				writeJSONResponse(w, http.StatusInternalServerError, response)
				return
			}
		}
	}

	// Check if response should include cards
	includeCards := shouldIncludeCards(r.URL.Query().Get("include"))
	if includeCards {
		cards, err := h.storage.ListDeckCards(r.Context(), deckID)
		if err != nil {
			h.logger.Error("Failed to get deck cards for response",
				slog.String("operation", "get_deck_cards"),
				slog.String("deck_id", deckIDStr),
				slog.Any("error", err))
			// Continue without cards in response rather than failing the whole request
		} else {
			totalCount := 0
			for _, cardWithQuantity := range cards {
				totalCount += cardWithQuantity.Quantity
			}

			cardCollection := &models.CardCollection{
				TotalCount:  totalCount,
				UniqueCount: len(cards),
				Items:       make([]models.CardWithQuantity, len(cards)),
			}

			// Convert from pointers to values for the response
			for i, cardPtr := range cards {
				cardCollection.Items[i] = *cardPtr
			}
			updatedDeck.Cards = cardCollection
		}
	}

	writeJSONResponse(w, http.StatusOK, updatedDeck)
}

// deleteDeck handles DELETE /decks/{id}
func (h *DecksHandler) deleteDeck(w http.ResponseWriter, r *http.Request, deckIDStr string) {
	// Parse deck ID
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid deck ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	err = h.storage.DeleteDeck(r.Context(), deckID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			response := ErrorResponse{
				Error:   errStrNotFound,
				Message: "Deck not found",
			}
			writeJSONResponse(w, http.StatusNotFound, response)
		} else {
			h.logger.Error("Failed to delete deck",
				slog.String("operation", "delete_deck"),
				slog.String("deck_id", deckIDStr),
				slog.Any("error", err))

			response := ErrorResponse{
				Error:   errStrInternal,
				Message: "Failed to delete deck",
			}
			writeJSONResponse(w, http.StatusInternalServerError, response)
		}
		return
	}

	// successful deletion
	w.WriteHeader(http.StatusNoContent)
}

// shouldIncludeCards parses the include parameter and returns true if "cards" is requested
func shouldIncludeCards(includeParam string) bool {
	for include := range strings.SplitSeq(includeParam, ",") {
		if strings.TrimSpace(include) == "cards" {
			return true
		}
	}
	return false
}
