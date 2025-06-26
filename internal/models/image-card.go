package models

import (
	"errors"
	"strings"
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
	FieldTypes: map[string]query.FieldType{
		"id":         query.FieldTypeUUID,
		"name":       query.FieldTypeString,
		"created_at": query.FieldTypeDateTime,
		"updated_at": query.FieldTypeDateTime,
	},
}

const (
	TypeImageCard   = "image-card"
	TypeImageCardID = 1
)

func (c *ImageCard) GetID() uuid.UUID         { return c.ID }
func (c *ImageCard) GetName() string          { return c.Name }
func (c *ImageCard) GetFrontImageURL() string { return c.FrontImageURL }
func (c *ImageCard) GetBackImageURL() string  { return c.BackImageURL }
func (c *ImageCard) GetCardType() string      { return TypeImageCard }

// Validate validates the ImageCard fields
func (c *ImageCard) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if len(c.Name) > 255 {
		return errors.New("name must be 255 characters or less")
	}
	if len(c.Description) > 1000 {
		return errors.New("description must be 1000 characters or less")
	}
	if len(c.FrontImageURL) > 500 {
		return errors.New("front_image_url must be 500 characters or less")
	}
	if len(c.BackImageURL) > 500 {
		return errors.New("back_image_url must be 500 characters or less")
	}
	if c.FrontImageURL != "" && !isValidURL(c.FrontImageURL) {
		return errors.New("front_image_url must be a valid URL")
	}
	if c.BackImageURL != "" && !isValidURL(c.BackImageURL) {
		return errors.New("back_image_url must be a valid URL")
	}
	return nil
}

// isValidURL performs basic URL validation
func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return true // Empty URLs are allowed
	}
	return strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")
}
