package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
)

// GameCard operations - stubs
func (m *MySQLStorage) ListGameCards(ctx context.Context, cardType string) ([]*deckdef.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) GetGameCard(ctx context.Context, id uuid.UUID) (*deckdef.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) CreateGameCard(ctx context.Context, card deckdef.GameCard) (*deckdef.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) UpdateGameCard(ctx context.Context, card deckdef.GameCard) (*deckdef.GameCard, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *MySQLStorage) DeleteGameCard(ctx context.Context, id uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
