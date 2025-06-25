package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
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

// ImageCardQueryConfig defines allowed filters and sorts for ImageCard queries
var ImageCardQueryConfig = query.QueryConfig{
	AllowedFilters: map[string]string{
		"id":         "uuid",
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	},
	AllowedSorts: map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
	},
}

const TypeImageCard = "imagecard"

func (c *ImageCard) GetID() uuid.UUID         { return c.ID }
func (c *ImageCard) GetName() string          { return c.Name }
func (c *ImageCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *ImageCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *ImageCard) GetCardType() string      { return TypeImageCard }
