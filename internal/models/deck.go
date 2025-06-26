package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// Deck represents the base deck structure shared across all games
type Deck struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	DeckType       *string         `json:"deck_type,omitempty"`
	SleeveImageURL *string         `json:"sleeve_image_url,omitempty"`
	Cards          []CardInterface `json:"cards,omitempty"` // TODO: Only populated when explicitly requested
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// DeckType constants for different types of decks
// matches db seeds for card types
const (
	DeckTypeStandard     = 1
	DeckTypePlayingCards = 2
)

// DeckQueryConfig defines allowed filters and sorts for Deck queries
var DeckQueryConfig = query.QueryConfig{
	AllowedFilters: map[string]string{
		"id":         "d.uuid",
		"name":       "d.name",
		"deck_type":  "dt.name",
		"created_at": "d.created_at",
		"updated_at": "d.updated_at",
	},
	AllowedSorts: map[string]string{
		"name":       "d.name",
		"deck_type":  "dt.name",
		"created_at": "d.created_at",
		"updated_at": "d.updated_at",
	},
	FieldTypes: map[string]query.FieldType{
		"id":         query.FieldTypeUUID,
		"name":       query.FieldTypeString,
		"deck_type":  query.FieldTypeString,
		"created_at": query.FieldTypeDateTime,
		"updated_at": query.FieldTypeDateTime,
	},
}
