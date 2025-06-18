package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PlayingCard represents a standard playing card
type PlayingCard struct {
	ID            uuid.UUID `json:"id"`
	Suite         string    `json:"suite"` // Hearts, Diamonds, Clubs, Spades
	Value         int       `json:"value"` // 1-13 (Ace through King)
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *PlayingCard) GetID() uuid.UUID         { return c.ID }
func (c *PlayingCard) GetName() string          { return fmt.Sprintf("%s of %s", c.getValueName(), c.Suite) }
func (c *PlayingCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *PlayingCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *PlayingCard) GetCardType() string      { return "playingcard" }

func (c *PlayingCard) getValueName() string {
	switch c.Value {
	case 1:
		return "Ace"
	case 11:
		return "Jack"
	case 12:
		return "Queen"
	case 13:
		return "King"
	default:
		return fmt.Sprintf("%d", c.Value)
	}
}
