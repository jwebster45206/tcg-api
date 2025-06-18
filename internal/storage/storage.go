package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

type Storage interface {
	// Deck operations
	ListDecks(ctx context.Context, ownerID *uuid.UUID) ([]*models.Deck, error)
	GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error)
	CreateDeck(ctx context.Context, deck *models.Deck) error
	UpdateDeck(ctx context.Context, id uuid.UUID, deck *models.Deck) error
	DeleteDeck(ctx context.Context, id uuid.UUID) error

	// Generic card operations - cardType determines which table/collection
	ListCards(ctx context.Context, cardType string) ([]models.CardInterface, error)
	GetCard(ctx context.Context, id uuid.UUID, cardType string) (models.CardInterface, error)
	CreateCard(ctx context.Context, card models.CardInterface) error
	UpdateCard(ctx context.Context, id uuid.UUID, card models.CardInterface) error
	DeleteCard(ctx context.Context, id uuid.UUID, cardType string) error

	// TODO: DeckState operations for future gameplay mechanics
	// CreateDeckState, GetDeckState, UpdateDeckState, DeleteDeckState
}
