package models

import (
	"time"

	"github.com/google/uuid"
)

// GameCard represents a TCG-specific card with game mechanics
type GameCard struct {
	ID            uuid.UUID
	Name          string    `json:"name"`
	Subtitle      string    `json:"subtitle"`
	Cost          int       `json:"cost"`
	Type          string    `json:"type"`
	Offense       int       `json:"offense"`
	Defense       int       `json:"defense"`
	Keywords      []string  `json:"keywords"`
	Colors        []string  `json:"colors"`
	IsResource    bool      `json:"is_resource"`
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const TypeGameCard = "game-card"

// Implement CardInterface
func (c *GameCard) GetID() uuid.UUID         { return c.ID }
func (c *GameCard) GetName() string          { return c.Name }
func (c *GameCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *GameCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *GameCard) GetCardType() string      { return TypeGameCard }
