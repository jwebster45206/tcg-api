package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// PlayingCard operations - stubs
func (m *MySQLStorage) ListPlayingCards(ctx context.Context) ([]*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetPlayingCard(ctx context.Context, id uuid.UUID) (*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeletePlayingCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
