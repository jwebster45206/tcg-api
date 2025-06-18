package models

import (
	"time"

	"github.com/google/uuid"
)

// ImageCard represents a simple card with just imagery and basic info
type ImageCard struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	FrontImageURL string    `json:"front_image_url"`
	BackImageURL  string    `json:"back_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (c *ImageCard) GetID() uuid.UUID         { return c.ID }
func (c *ImageCard) GetName() string          { return c.Name }
func (c *ImageCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *ImageCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *ImageCard) GetCardType() string      { return "imagecard" }
