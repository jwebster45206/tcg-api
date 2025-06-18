package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// Different storage interfaces for different domains
type Storage interface {
	ListCards(ctx context.Context) ([]*models.GameCard, error)
	GetCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error)
	CreateCard(ctx context.Context, card *models.GameCard) error
	UpdateCard(ctx context.Context, id uuid.UUID, card *models.GameCard) error
	DeleteCard(ctx context.Context, id uuid.UUID) error
}
