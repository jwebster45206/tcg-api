package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// GameCard operations - stubs
func (m *MySQLStorage) ListGameCards(ctx context.Context, cardType string) ([]*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetGameCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteGameCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
