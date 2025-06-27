package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/query"
)

type Storage interface {
	// Health check
	Ping(ctx context.Context) error

	// Deck operations
	ListDecks(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.Deck, error)
	GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error)
	CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error)
	UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error)
	DeleteDeck(ctx context.Context, id uuid.UUID) error
	ListDeckCards(ctx context.Context, deckID uuid.UUID) ([]*models.CardWithQuantity, error)
	SetDeckCards(ctx context.Context, deckID uuid.UUID, cards []models.DeckCardInput) error

	// ImageCard operations
	ListImageCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.ImageCard, error)
	GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error)
	CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error)
	UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error)
	DeleteImageCard(ctx context.Context, id uuid.UUID) error

	// PlayingCard operations
	ListPlayingCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.PlayingCard, error)
	GetPlayingCard(ctx context.Context, id uuid.UUID) (*models.PlayingCard, error)
	CreatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error)
	UpdatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error)
	DeletePlayingCard(ctx context.Context, id uuid.UUID) error

	// GameCard operations
	ListGameCards(ctx context.Context, cardType string) ([]*models.GameCard, error)
	GetGameCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error)
	CreateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error)
	UpdateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error)
	DeleteGameCard(ctx context.Context, id uuid.UUID) error

	// TODO: DeckState operations for future gameplay mechanics
	// CreateDeckState, GetDeckState, UpdateDeckState, DeleteDeckState
}

// Helper function to safely dereference nullable string pointers
// Returns empty string if the pointer is nil
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to safely dereference nullable int pointers
// Returns 0 if the pointer is nil
func safeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
