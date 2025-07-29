package deckdef

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// PlayingCard represents a standard playing card
type PlayingCard struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	Suit          string    `json:"suit"`    // Hearts, Diamonds, Clubs, Spades
	Ranking       int       `json:"ranking"` // 1-13 (Ace through King)
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PlayingCardQueryConfig defines allowed filters and sorts for PlayingCard queries
var PlayingCardQueryConfig = query.QueryConfig{
	AllowedFilters: map[string]string{
		"id":         "c.uuid",
		"suit":       "pc.suit",
		"ranking":    "pc.ranking",
		"created_at": "c.created_at",
		"updated_at": "c.updated_at",
	},
	AllowedSorts: map[string]string{
		"suit":       "pc.suit",
		"ranking":    "pc.ranking",
		"created_at": "c.created_at",
		"updated_at": "c.updated_at",
	},
	FieldTypes: map[string]query.FieldType{
		"id":         query.FieldTypeUUID,
		"suit":       query.FieldTypeString,
		"ranking":    query.FieldTypeInt,
		"created_at": query.FieldTypeDateTime,
		"updated_at": query.FieldTypeDateTime,
	},
}

const (
	TypePlayingCard   = "playing-card"
	TypePlayingCardID = 2

	SuitHearts   = "hearts"
	SuitDiamonds = "diamonds"
	SuitClubs    = "clubs"
	SuitSpades   = "spades"

	ColorRed   = "red"
	ColorBlack = "black"

	ValueAce   = "ace"
	ValueJack  = "jack"
	ValueQueen = "queen"
	ValueKing  = "king"
)

func (c *PlayingCard) GetID() uuid.UUID         { return c.ID }
func (c *PlayingCard) GetName() string          { return fmt.Sprintf("%s of %s", c.getValueName(), c.Suit) }
func (c *PlayingCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *PlayingCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *PlayingCard) GetCardType() string      { return TypePlayingCard }

func (c *PlayingCard) getValueName() string {
	switch c.Ranking {
	case 1:
		return ValueAce
	case 11:
		return ValueJack
	case 12:
		return ValueQueen
	case 13:
		return ValueKing
	default:
		return fmt.Sprintf("%d", c.Ranking)
	}
}

func (c *PlayingCard) GetColor() string {
	switch c.Suit {
	case SuitHearts, SuitDiamonds:
		return ColorRed
	case SuitClubs, SuitSpades:
		return ColorBlack
	default:
		return "unknown"
	}
}

func (c *PlayingCard) Validate() error {
	if c.Ranking < 1 || c.Ranking > 13 {
		return fmt.Errorf("ranking must be between 1 and 13")
	}
	if c.Suit != SuitHearts &&
		c.Suit != SuitDiamonds &&
		c.Suit != SuitClubs &&
		c.Suit != SuitSpades {
		return fmt.Errorf("invalid suit: %s", c.Suit)
	}
	return nil
}
