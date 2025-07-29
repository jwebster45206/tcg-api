package deckdef

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GameCard represents a generic TCG card with common game mechanics
type GameCard struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ManaCost      int       `json:"mana_cost"`
	CardType      string    `json:"card_type"`
	Attack        int       `json:"attack"`
	Health        int       `json:"health"`
	Keywords      []string  `json:"keywords"`
	Colors        []string  `json:"colors"`
	Rarity        string    `json:"rarity"`
	SetCode       string    `json:"set_code,omitempty"`
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const (
	TypeGameCard = "game-card"
	// TypeGameCardID = 3
)

// Implement CardInterface
func (c *GameCard) GetID() uuid.UUID         { return c.ID }
func (c *GameCard) GetName() string          { return c.Name }
func (c *GameCard) GetDescription() string   { return c.Description }
func (c *GameCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *GameCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *GameCard) GetCardType() string      { return TypeGameCard }

// Validate validates the GameCard fields
func (c *GameCard) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if len(c.Name) > 255 {
		return errors.New("name must be 255 characters or less")
	}
	if len(c.Description) > 1000 {
		return errors.New("description must be 1000 characters or less")
	}
	if c.ManaCost < 0 {
		return errors.New("mana_cost cannot be negative")
	}
	if c.Attack < 0 {
		return errors.New("attack cannot be negative")
	}
	if c.Health < 0 {
		return errors.New("health cannot be negative")
	}
	if c.Rarity != "" && !isValidRarity(c.Rarity) {
		return errors.New("rarity must be one of: common, uncommon, rare, mythic")
	}
	if len(c.SetCode) > 10 {
		return errors.New("set_code must be 10 characters or less")
	}
	return nil
}

// isValidRarity checks if the rarity is one of the standard TCG rarities
func isValidRarity(rarity string) bool {
	validRarities := []string{"common", "uncommon", "rare", "mythic", "legendary"}
	for _, validRarity := range validRarities {
		if rarity == validRarity {
			return true
		}
	}
	return false
}
