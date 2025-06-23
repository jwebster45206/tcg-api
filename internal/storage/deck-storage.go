package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// Deck operations - stubs
func (m *MySQLStorage) ListDecks(ctx context.Context, ownerID *uuid.UUID) ([]*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteDeck(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
