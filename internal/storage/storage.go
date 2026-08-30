package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/query"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
)

type Storage interface {
	// Health check
	Ping(ctx context.Context) error

	// Deck operations
	ListDecks(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*deckdef.Deck, error)
	GetDeck(ctx context.Context, id uuid.UUID) (*deckdef.Deck, error)
	CreateDeck(ctx context.Context, deck deckdef.Deck) (*deckdef.Deck, error)
	UpdateDeck(ctx context.Context, deck deckdef.Deck) (*deckdef.Deck, error)
	DeleteDeck(ctx context.Context, id uuid.UUID) error
	ListDeckCards(ctx context.Context, deckID uuid.UUID) ([]*deckdef.CardWithQuantity, error)
	SetDeckCards(ctx context.Context, deckID uuid.UUID, cards []deckdef.CardInputWithQuantity) error

	// ImageCard operations
	ListImageCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*deckdef.ImageCard, error)
	GetImageCard(ctx context.Context, id uuid.UUID) (*deckdef.ImageCard, error)
	CreateImageCard(ctx context.Context, imageCard deckdef.ImageCard) (*deckdef.ImageCard, error)
	UpdateImageCard(ctx context.Context, imageCard deckdef.ImageCard) (*deckdef.ImageCard, error)
	DeleteImageCard(ctx context.Context, id uuid.UUID) error

	// PlayingCard operations
	ListPlayingCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*deckdef.PlayingCard, error)
	GetPlayingCard(ctx context.Context, id uuid.UUID) (*deckdef.PlayingCard, error)
	CreatePlayingCard(ctx context.Context, card deckdef.PlayingCard) (*deckdef.PlayingCard, error)
	UpdatePlayingCard(ctx context.Context, card deckdef.PlayingCard) (*deckdef.PlayingCard, error)
	DeletePlayingCard(ctx context.Context, id uuid.UUID) error
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
