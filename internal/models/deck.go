package models

import (
	"time"

	"github.com/google/uuid"
)

// Deck represents the base deck structure shared across all games
type Deck struct {
	ID             uuid.UUID   `json:"id"`
	Name           string      `json:"name"`
	OwnerID        *uuid.UUID  `json:"owner_id,omitempty"`
	SleeveImageURL *string     `json:"sleeve_image_url,omitempty"`
	BackImageURL   *string     `json:"back_image_url,omitempty"`
	Cards          []uuid.UUID `json:"cards"` // Array of card identifiers
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
