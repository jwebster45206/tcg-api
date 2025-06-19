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
	ListGameCards(ctx context.Context, cardType string) ([]models.GameCard, error)
	GetGameCard(ctx context.Context, id uuid.UUID, cardType string) (models.GameCard, error)
	CreateGameCard(ctx context.Context, card models.GameCard) (models.GameCard, error)
	UpdateGameCard(ctx context.Context, id uuid.UUID, card models.GameCard) error
	DeleteGameCard(ctx context.Context, id uuid.UUID) error

	// TODO: DeckState operations for future gameplay mechanics
	// CreateDeckState, GetDeckState, UpdateDeckState, DeleteDeckState
}
