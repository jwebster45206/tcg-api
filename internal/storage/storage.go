package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// FilterOperator represents SQL comparison operators
type FilterOperator string

const (
	OpEqual        FilterOperator = "="  // Auto converts to "IN" for arrays, "IS NULL" for nil
	OpNotEqual     FilterOperator = "!=" // Auto converts to "NOT IN" for arrays, "IS NOT NULL" for nil
	OpGreaterThan  FilterOperator = ">"
	OpGreaterEqual FilterOperator = ">="
	OpLessThan     FilterOperator = "<"
	OpLessEqual    FilterOperator = "<="
	OpLike         FilterOperator = "LIKE"
	OpNotLike      FilterOperator = "NOT LIKE"
)

// Filter represents a single filter condition
type Filter struct {
	Column   string
	Operator FilterOperator
	Value    interface{} // nil, single value, or slice
}

// SortOption represents a field to sort by and its direction
type SortOption struct {
	Field string
	Desc  bool // true for DESC, false for ASC
}

type Storage interface {
	// Health check
	Ping(ctx context.Context) error

	// Deck operations
	ListDecks(ctx context.Context, ownerID *uuid.UUID) ([]*models.Deck, error)
	GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error)
	CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error)
	UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error)
	DeleteDeck(ctx context.Context, id uuid.UUID) error

	// ImageCard operations
	ListImageCards(ctx context.Context, filters []Filter, sorts []SortOption, pageSize int, pageNum int) ([]*models.ImageCard, error)
	GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error)
	CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error)
	UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error)
	DeleteImageCard(ctx context.Context, id uuid.UUID) error

	// PlayingCard operations
	ListPlayingCards(ctx context.Context) ([]*models.PlayingCard, error)
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
