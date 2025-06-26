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
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	sorts, err := ParseSorts(r, models.DeckQueryConfig)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	pageSize, pageNum, err := ParsePagination(r)
	if err != nil {
		h.handleError(w, err, http.StatusBadRequest)
		return
	}
	decks, err := h.storage.ListDecks(r.Context(), filters, sorts, pageSize, pageNum)
	if err != nil {
		h.logger.Error("Failed to list decks", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
			http.Error(w, "Including cards is only allowed when filtering by 'id'", http.StatusBadRequest)
			return
		}
	}

	// If cards are requested, fetch and populate them for each deck
	if includeCards {
		for _, deck := range decks {
			cards, err := h.storage.ListDeckCards(r.Context(), deck.ID)
			if err != nil {
				h.logger.Error("Failed to get deck cards", "error", err, "deck_id", deck.ID)
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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(decks); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// getDeck handles GET /decks/{id}
func (h *DecksHandler) getDeck(w http.ResponseWriter, r *http.Request, deckIDStr string) {
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		http.Error(w, "Invalid deck ID format", http.StatusBadRequest)
		return
	}

	deck, err := h.storage.GetDeck(r.Context(), deckID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Deck not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to get deck", "error", err, "deck_id", deckID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// If cards are requested, fetch and populate them
	includeCards := shouldIncludeCards(r.URL.Query().Get("include"))
	if includeCards {
		cards, err := h.storage.ListDeckCards(r.Context(), deckID)
		if err != nil {
			h.logger.Error("Failed to get deck cards", "error", err, "deck_id", deckID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(deck); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// createDeck handles POST /decks
func (h *DecksHandler) createDeck(w http.ResponseWriter, r *http.Request) {
	var deck models.Deck
	if err := json.NewDecoder(r.Body).Decode(&deck); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdDeck, err := h.storage.CreateDeck(r.Context(), deck)
	if err != nil {
		h.logger.Error("Failed to create deck", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdDeck); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// updateDeck handles PATCH /decks/{id}
func (h *DecksHandler) updateDeck(w http.ResponseWriter, r *http.Request, deckIDStr string) {
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		http.Error(w, "Invalid deck ID format", http.StatusBadRequest)
		return
	}

	// Parse request body
	var deck models.Deck
	if err := json.NewDecoder(r.Body).Decode(&deck); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deck.ID = deckID
	updatedDeck, err := h.storage.UpdateDeck(r.Context(), deck)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Deck not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to update deck", "error", err, "deck_id", deckID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedDeck); err != nil {
		h.logger.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// deleteDeck handles DELETE /decks/{id}
func (h *DecksHandler) deleteDeck(w http.ResponseWriter, r *http.Request, deckIDStr string) {
	// Parse deck ID
	deckID, err := uuid.Parse(deckIDStr)
	if err != nil {
		http.Error(w, "Invalid deck ID format", http.StatusBadRequest)
		return
	}

	err = h.storage.DeleteDeck(r.Context(), deckID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Deck not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to delete deck", "error", err, "deck_id", deckID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// successful deletion
	w.WriteHeader(http.StatusNoContent)
}

// handleError handles validation errors and other structured errors
func (h *DecksHandler) handleError(w http.ResponseWriter, err error, defaultStatus int) {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if encErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "validation_error",
			"field":   validationErr.Field,
			"message": validationErr.Message,
		}); encErr != nil {
			h.logger.Error("Failed to encode validation error response", "error", encErr)
		}
		return
	}

	http.Error(w, err.Error(), defaultStatus)
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
