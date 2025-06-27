package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// DeckType constants for different types of decks
// matches db seeds for card types
const (
	DeckTypeStandard    = "standard"
	DeckTypePlayingCard = "playing-card"
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

// Deck represents the base deck structure shared across all games
type Deck struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	DeckType       string          `json:"deck_type"`
	SleeveImageURL *string         `json:"sleeve_image_url,omitempty"`
	Cards          *CardCollection `json:"cards,omitempty"` // Only when ?include=cards
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// CardCollection is a representation of cards in a deck
// with de-duplication and quantity tracking.
type CardCollection struct {
	TotalCount  int                `json:"total_count"`  // Sum of all quantities (54 for standard deck)
	UniqueCount int                `json:"unique_count"` // Number of different cards (54 for standard deck)
	Items       []CardWithQuantity `json:"items"`
}

type CardWithQuantity struct {
	Card     CardInterface `json:"card"`
	Quantity int           `json:"quantity"`
}

// OrderedCardCollection is a representation of cards in a deck
// with their positions for ordered decks (like playing cards).
// Duplication is allowed.
// type OrderedCardCollection struct {
// 	TotalCount int              `json:"total_count"`
// 	Items      []PositionedCard `json:"items"`
// }

// type PositionedCard struct {
// 	Card     CardInterface `json:"card"`
// 	Position int           `json:"position"`
// }

// DeckCardInput represents a simplified card input for create/update operations
// Clients only need to provide card ID and quantity
type DeckCardInput struct {
	CardID   uuid.UUID `json:"card_id"`
	Quantity int       `json:"quantity"`
}

// DeckInput represents the input structure for deck creation and updates
// Supports optional card management during deck operations
type DeckInput struct {
	Name           string               `json:"name"`
	DeckType       string               `json:"deck_type"`
	SleeveImageURL *string              `json:"sleeve_image_url,omitempty"`
	Cards          *CardCollectionInput `json:"cards,omitempty"` // Optional cards for create/update
}

// CardCollectionInput is used for create/update operations with simplified card inputs
type CardCollectionInput struct {
	Items []DeckCardInput `json:"items"`
}

// Validate validates the DeckInput and returns a ValidationError if invalid
func (d *DeckInput) Validate() error {
	// Validate required fields
	if d.Name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "Deck name is required",
		}
	}

	// Validate card quantities if cards are provided
	if d.Cards != nil {
		for i, cardInput := range d.Cards.Items {
			if cardInput.CardID == uuid.Nil {
				return &ValidationError{
					Field:   fmt.Sprintf("cards.items[%d].card_id", i),
					Message: "Card ID is required",
				}
			}
			if cardInput.Quantity <= 0 {
				return &ValidationError{
					Field:   fmt.Sprintf("cards.items[%d].quantity", i),
					Message: "Card quantity must be positive",
				}
			}
		}
	}

	return nil
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
