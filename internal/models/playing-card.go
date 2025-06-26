package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
)

// PlayingCard represents a standard playing card
type PlayingCard struct {
	ID            uuid.UUID `json:"id"`
	Suite         string    `json:"suite"`   // Hearts, Diamonds, Clubs, Spades
	Ranking       int       `json:"ranking"` // 1-13 (Ace through King)
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PlayingCardQueryConfig defines allowed filters and sorts for PlayingCard queries
var PlayingCardQueryConfig = query.QueryConfig{
	AllowedFilters: map[string]string{
		"id":         "c.uuid",
		"suite":      "pc.suit",
		"ranking":    "pc.ranking",
		"created_at": "c.created_at",
		"updated_at": "c.updated_at",
	},
	AllowedSorts: map[string]string{
		"suite":      "pc.suit",
		"ranking":    "pc.ranking",
		"created_at": "c.created_at",
		"updated_at": "c.updated_at",
	},
	FieldTypes: map[string]query.FieldType{
		"id":         query.FieldTypeUUID,
		"suite":      query.FieldTypeString,
		"ranking":    query.FieldTypeInt,
		"created_at": query.FieldTypeDateTime,
		"updated_at": query.FieldTypeDateTime,
	},
}

const (
	TypePlayingCard   = "playing-card"
	TypePlayingCardID = 2

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
	if c.Ranking < 1 || c.Ranking > 13 {
		return fmt.Errorf("ranking must be between 1 and 13")
	}
	if c.Suite != SuiteHearts &&
		c.Suite != SuiteDiamonds &&
		c.Suite != SuiteClubs &&
		c.Suite != SuiteSpades {
		return fmt.Errorf("invalid suite: %s", c.Suite)
	}
	return nil
}
