package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// DeckType constants for different types of decks
// matches db seeds for deck types
const (
	DeckTypeStandard    = "standard"
	DeckTypePlayingCard = "playing-card"
)

// DeckQueryConfig defines allowed filters and sorts for Deck queries
var DeckQueryConfig = query.QueryConfig{
	AllowedFilters: map[string]string{
		"id":         "d.uuid",
		"name":       "d.name",
		"type":       "dt.name",
		"created_at": "d.created_at",
		"updated_at": "d.updated_at",
	},
	AllowedSorts: map[string]string{
		"name":       "d.name",
		"type":       "dt.name",
		"created_at": "d.created_at",
		"updated_at": "d.updated_at",
	},
	FieldTypes: map[string]query.FieldType{
		"id":         query.FieldTypeUUID,
		"name":       query.FieldTypeString,
		"type":       query.FieldTypeString,
		"created_at": query.FieldTypeDateTime,
		"updated_at": query.FieldTypeDateTime,
	},
}

// Deck represents the base deck structure shared across all games
type Deck struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
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

// CardInput represents a card reference for input operations.
// It is a generic type that can be used for any card type for input purposes.
type CardInput struct {
	ID uuid.UUID `json:"id"`
}

// CardInputWithQuantity represents a simplified card input for create/update operations
// It is a generic version of CardWithQuantity, used for input purposes.
type CardInputWithQuantity struct {
	Card     CardInput `json:"card"`
	Quantity int       `json:"quantity"`
}

// CardCollectionInput is a simplified version of CardCollection.
// As its name suggests, is used for input purposes.
type CardCollectionInput struct {
	Items []CardInputWithQuantity `json:"items"`
}

// DeckInput represents the input structure for deck creation and updates
// Supports optional card management during deck operations
type DeckInput struct {
	Name           string               `json:"name"`
	Type           string               `json:"type"`
	SleeveImageURL *string              `json:"sleeve_image_url,omitempty"`
	Cards          *CardCollectionInput `json:"cards,omitempty"` // Optional cards for create/update
}

func (d *DeckInput) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("deck name is required")
	}
	if d.Cards != nil {
		for i, cardInput := range d.Cards.Items {
			if cardInput.Card.ID == uuid.Nil {
				return fmt.Errorf("card ID is required for cards.items[%d]", i)
			}
			if cardInput.Quantity <= 0 {
				return fmt.Errorf("card quantity must be positive for cards.items[%d]", i)
			}
		}
	}
	return nil
}
