package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// ImageCard operations - stubs
func (m *MySQLStorage) ListImageCards(ctx context.Context) ([]*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteImageCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
