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

const (
	TypePlayingCard = "playing-card"

	SuiteHearts   = "hearts"
	SuiteDiamonds = "diamonds"
	SuiteClubs    = "clubs"
	SuiteSpades   = "spades"

	ColorRed   = "red"
	ColorBlack = "black"

	ValueAce   = "ace"
	ValueJack  = "jack"
	ValueQueen = "queen"
	ValueKing  = "king"
)

func (c *PlayingCard) GetID() uuid.UUID         { return c.ID }
func (c *PlayingCard) GetName() string          { return fmt.Sprintf("%s of %s", c.getValueName(), c.Suite) }
func (c *PlayingCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *PlayingCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *PlayingCard) GetCardType() string      { return TypePlayingCard }

func (c *PlayingCard) getValueName() string {
	switch c.Value {
	case 1:
		return ValueAce
	case 11:
		return ValueJack
	case 12:
		return ValueQueen
	case 13:
		return ValueKing
	default:
		return fmt.Sprintf("%d", c.Value)
	}
}

func (c *PlayingCard) GetColor() string {
	switch c.Suite {
	case SuiteHearts, SuiteDiamonds:
		return ColorRed
	case SuiteClubs, SuiteSpades:
		return ColorBlack
	default:
		return "unknown"
	}
}

func (c *PlayingCard) Validate() error {
	if c.Value < 1 || c.Value > 13 {
		return fmt.Errorf("value must be between 1 and 13")
	}
	if c.Suite != SuiteHearts && c.Suite != SuiteDiamonds && c.Suite != SuiteClubs && c.Suite != SuiteSpades {
		return fmt.Errorf("invalid suite: %s", c.Suite)
	}
	return nil
}
