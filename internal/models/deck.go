package models

import (
	"time"

	"github.com/google/uuid"
)

// Deck represents the base deck structure shared across all games
type Deck struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	DeckTypeID     int             `json:"deck_type_id"`
	SleeveImageURL *string         `json:"sleeve_image_url,omitempty"`
	Cards          []CardInterface `json:"cards,omitempty"` // Cards in this deck with quantities
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// DeckType constants for different types of decks
// matches db seeds for card types
const (
	DeckTypeStandard     = 1
	DeckTypePlayingCards = 2
)
