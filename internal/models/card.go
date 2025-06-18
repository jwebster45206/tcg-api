package models

import "github.com/google/uuid"

// CardInterface defines the contract that all card types must implement
type CardInterface interface {
	GetID() uuid.UUID
	GetName() string
	GetFrontImageURL() string
	GetBackImageURL() string
	GetCardType() string // Used for routing to correct storage/handlers
}
