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
	CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error)
	UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error)
	DeleteDeck(ctx context.Context, id uuid.UUID) error

	// ImageCard operations
	ListImageCards(ctx context.Context) ([]*models.ImageCard, error)
	GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error)
	CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error)
	UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error)
	DeleteImageCard(ctx context.Context, id uuid.UUID) error

	// GameCard operations
	ListGameCards(ctx context.Context, cardType string) ([]*models.GameCard, error)
	GetGameCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error)
	CreateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error)
	UpdateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error)
	DeleteGameCard(ctx context.Context, id uuid.UUID) error

	// TODO: DeckState operations for future gameplay mechanics
	// CreateDeckState, GetDeckState, UpdateDeckState, DeleteDeckState
}
