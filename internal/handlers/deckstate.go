package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
	"github.com/jwebster45206/tcg-api/internal/state"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// ContextKeyDeckStateLock is the key for storing deck state locks in context
	ContextKeyDeckStateLock ContextKey = "deckstate_lock"
)

// DeckStateHandler handles deck state operations
type DeckStateHandler struct {
	storage      storage.Storage        // MySQL for deck definitions
	stateStorage state.DeckStateStorage // Redis for deck states
	logger       *slog.Logger
}

// NewDeckStateHandler creates a new DeckStateHandler with dependencies
func NewDeckStateHandler(storage storage.Storage, stateStorage state.DeckStateStorage, logger *slog.Logger) *DeckStateHandler {
	return &DeckStateHandler{
		storage:      storage,
		stateStorage: stateStorage,
		logger:       logger,
	}
}

// withLock is a middleware that acquires a pessimistic lock on the deck state before proceeding.
// This prevents race conditions when multiple players or systems try to modify the same deck state.
// The lock is automatically released when the handler completes.
// Returns HTTP 409 Conflict if the deck state is already locked by another process.
func (h *DeckStateHandler) withLock(next func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, stateID string) {
		// Acquire lock with 30 second timeout
		lock, err := h.stateStorage.LockDeckState(r.Context(), stateID, 30*time.Second)
		if err != nil {
			// Check if it's a lock error (resource locked) or other error
			if lockErr, ok := err.(*state.LockError); ok {
				h.logger.Warn("Deck state locked",
					slog.String("operation", "acquire_lock"),
					slog.String("deck_state_id", stateID),
					slog.String("error", lockErr.Message))
				response := ErrorResponse{
					Error:   "resource_locked",
					Message: "Deck state is being modified by another player. Please try again.",
				}
				writeJSONResponse(w, http.StatusConflict, response)
				return
			}

			// Other lock acquisition error
			h.logger.Error("Failed to acquire deck state lock",
				slog.String("operation", "acquire_lock"),
				slog.String("deck_state_id", stateID),
				slog.Any("error", err))
			response := ErrorResponse{
				Error:   errStrInternal,
				Message: "Failed to acquire resource lock",
			}
			writeJSONResponse(w, http.StatusInternalServerError, response)
			return
		}
		defer func() {
			if err := lock.Release(); err != nil {
				h.logger.Error("Failed to release deck state lock",
					slog.String("operation", "release_lock"),
					slog.String("deck_state_id", stateID),
					slog.String("lock_owner", lock.OwnerID()),
					slog.Any("error", err))
			}
		}()

		// Store lock in context for handlers to access if needed
		ctx := context.WithValue(r.Context(), ContextKeyDeckStateLock, lock)
		next(w, r.WithContext(ctx), stateID)
	}
}

// extractStateIDAndAction extracts the state ID and action from a path like "/{stateID}/actions/{action}"
func (h *DeckStateHandler) extractStateIDAndAction(path string) (string, string) {
	// path format: "/{stateID}/actions/{action}"
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 && parts[1] == "actions" {
		return parts[0], strings.Join(parts[2:], "/")
	}
	return "", ""
}

// handleActionWithStateID handles actions with the state ID already extracted and lock acquired
func (h *DeckStateHandler) handleActionWithStateID(w http.ResponseWriter, r *http.Request, stateID string) {
	// Reconstruct the original path for the existing handleAction method
	_, actionPath := h.extractStateIDAndAction(r.URL.Path[len("/v1/deckstates"):])
	fullPath := "/" + stateID + "/actions/" + actionPath
	h.handleAction(w, r, fullPath)
}

// deleteDeckStateWithLock handles deletion with lock already acquired
func (h *DeckStateHandler) deleteDeckStateWithLock(w http.ResponseWriter, r *http.Request, stateID string) {
	h.deleteDeckState(w, r, stateID)
}

func (h *DeckStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/deckstates")

	switch r.Method {
	case http.MethodGet:
		if path != "" && path != "/" {
			// Check if it's an action endpoint
			if strings.Contains(path, "/actions/") {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			stateID := strings.Trim(path, "/")
			h.getDeckState(w, r, stateID)
		} else {
			http.Error(w, "Deck state ID required", http.StatusBadRequest)
		}

	case http.MethodPost:
		if path == "" || path == "/" {
			h.createDeckState(w, r)
		} else if strings.Contains(path, "/actions/") {
			// Extract stateID from path and apply locking middleware
			stateID, _ := h.extractStateIDAndAction(path)
			if stateID == "" {
				http.Error(w, "Invalid deck state ID", http.StatusBadRequest)
				return
			}
			h.withLock(h.handleActionWithStateID)(w, r, stateID)
		} else {
			http.Error(w, "Method not allowed for this path", http.StatusMethodNotAllowed)
		}

	case http.MethodDelete:
		if path != "" && path != "/" {
			// Check if it's an action endpoint
			if strings.Contains(path, "/actions/") {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			stateID := strings.Trim(path, "/")
			h.withLock(h.deleteDeckStateWithLock)(w, r, stateID)
		} else {
			http.Error(w, "Deck state ID required for deletion", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// CreateDeckStateRequest represents the request to create a new deck state
type CreateDeckStateRequest struct {
	DeckID      uuid.UUID `json:"deck_id"`
	PlayerCount int       `json:"player_count"`
}

// createDeckState expands a deck into an immutable deck definition,
// and creates a new deck state in Redis.
func (h *DeckStateHandler) createDeckState(w http.ResponseWriter, r *http.Request) {
	var req CreateDeckStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid JSON in request body",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if req.DeckID == uuid.Nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "deck_id is required",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	if req.PlayerCount < 1 || req.PlayerCount > 8 {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "player_count must be between 1 and 8",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()

	deck, err := h.storage.GetDeck(ctx, req.DeckID)
	if err != nil {
		h.logger.Error("Failed to get deck for state creation",
			slog.String("operation", "get_deck"),
			slog.String("deck_id", req.DeckID.String()),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrNotFound,
			Message: "Deck not found",
		}
		writeJSONResponse(w, http.StatusNotFound, response)
		return
	}

	err = includeDeckCards(ctx, h.storage, deck)
	if err != nil {
		h.logger.Error("Failed to include deck cards",
			slog.String("operation", "include_deck_cards"),
			slog.String("deck_id", req.DeckID.String()),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to include deck cards",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Create new deck state using the deck definition
	deckState := deckstate.NewDeckState(*deck, req.PlayerCount)

	// Save the deck state to Redis
	if err := h.stateStorage.SaveDeckState(ctx, deckState.ID, deckState); err != nil {
		h.logger.Error("Failed to save deck state",
			slog.String("operation", "save_deck_state"),
			slog.String("deck_state_id", deckState.ID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to create deck state",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	h.logger.Info("Created new deck state",
		slog.String("operation", "create_deck_state"),
		slog.String("deck_state_id", deckState.ID),
		slog.String("deck_id", req.DeckID.String()),
		slog.Int("player_count", req.PlayerCount))

	writeJSONResponse(w, http.StatusCreated, deckState)
}

func (h *DeckStateHandler) getDeckState(w http.ResponseWriter, r *http.Request, stateID string) {
	// Validate UUID format
	_, err := uuid.Parse(stateID)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid deck state ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()
	deckState, err := h.stateStorage.GetDeckState(ctx, stateID)
	if err != nil {
		h.logger.Error("Failed to get deck state",
			slog.String("operation", "get_deck_state"),
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

	writeJSONResponse(w, http.StatusOK, deckState)
}

func (h *DeckStateHandler) deleteDeckState(w http.ResponseWriter, r *http.Request, stateID string) {
	// Validate UUID format
	_, err := uuid.Parse(stateID)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid deck state ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	ctx := r.Context()

	// Check if deck state exists first
	deckState, err := h.stateStorage.GetDeckState(ctx, stateID)
	if err != nil {
		h.logger.Error("Failed to check deck state existence",
			slog.String("operation", "get_deck_state_for_delete"),
			slog.String("deck_state_id", stateID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to check deck state",
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

	if err := h.stateStorage.DeleteDeckState(ctx, stateID); err != nil {
		h.logger.Error("Failed to delete deck state",
			slog.String("operation", "delete_deck_state"),
			slog.String("deck_state_id", stateID),
			slog.Any("error", err))
		response := ErrorResponse{
			Error:   errStrInternal,
			Message: "Failed to delete deck state",
		}
		writeJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	h.logger.Info("Deleted deck state",
		slog.String("operation", "delete_deck_state"),
		slog.String("deck_state_id", stateID))

	w.WriteHeader(http.StatusNoContent)
}

// handleAction routes action requests to the appropriate handler
func (h *DeckStateHandler) handleAction(w http.ResponseWriter, r *http.Request, path string) {
	// Parse path: /{id}/actions/{actionName}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[1] != "actions" {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid action path format. Expected /{id}/actions/{actionName}",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	stateID := parts[0]
	actionName := parts[2]

	// Validate UUID format
	_, err := uuid.Parse(stateID)
	if err != nil {
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Invalid deck state ID format",
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Route to specific action handler
	switch actionName {
	case "add-zone":
		h.handleAddZone(w, r, stateID)
	case "remove-zone":
		h.handleRemoveZone(w, r, stateID)
	case "sort-zone":
		h.handleSortZone(w, r, stateID)
	default:
		response := ErrorResponse{
			Error:   errStrBadRequest,
			Message: "Unsupported action: " + actionName,
		}
		writeJSONResponse(w, http.StatusBadRequest, response)
	}
}

// handleRemoveZone removes a zone from a deck state
func (h *DeckStateHandler) handleRemoveZone(w http.ResponseWriter, r *http.Request, stateID string) {
	h.logger.Info("Remove zone action requested",
		slog.String("operation", "remove_zone"),
		slog.String("deck_state_id", stateID))

	response := ErrorResponse{
		Error:   errStrNotImplemented,
		Message: "Remove zone action not implemented yet",
	}
	writeJSONResponse(w, http.StatusNotImplemented, response)
}
