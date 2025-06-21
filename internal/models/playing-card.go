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
	Value         int       `json:"value"` // 0-13 (Ace through King)
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const typePlayingCard = "playing-card"
const suiteHearts = "hearts"
const suiteDiamonds = "diamonds"
const suiteClubs = "clubs"
const suiteSpades = "spades"
const colorRed = "red"
const colorBlack = "black"
const valueAce = "ace"
const valueJack = "jack"
const valueQueen = "queen"
const valueKing = "king"

func (c *PlayingCard) GetID() uuid.UUID         { return c.ID }
func (c *PlayingCard) GetName() string          { return fmt.Sprintf("%s of %s", c.getValueName(), c.Suite) }
func (c *PlayingCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *PlayingCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *PlayingCard) GetCardType() string      { return typePlayingCard }

func (c *PlayingCard) getValueName() string {
	switch c.Value {
	case 1:
		return valueAce
	case 11:
		return valueJack
	case 12:
		return valueQueen
	case 13:
		return valueKing
	default:
		return fmt.Sprintf("%d", c.Value)
	}
}

func (c *PlayingCard) GetColor() string {
	switch c.Suite {
	case suiteHearts, suiteDiamonds:
		return colorRed
	case suiteClubs, suiteSpades:
		return colorBlack
	default:
		return "unknown"
	}
}

func (c *PlayingCard) Validate() error {
	if c.Value < 1 || c.Value > 13 {
		return fmt.Errorf("value must be between 1 and 13")
	}
	if c.Suite != suiteHearts && c.Suite != suiteDiamonds && c.Suite != suiteClubs && c.Suite != suiteSpades {
		return fmt.Errorf("invalid suite: %s", c.Suite)
	}
	return nil
}
