package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
)

// MockStorage implements CardStorage interface for testing and development
type MockStorage struct {
	mu    sync.RWMutex
	cards map[uuid.UUID]*models.GameCard
}

// NewMockStorage creates a new MockStorage instance with some sample data
func NewMockStorage() Storage {
	storage := &MockStorage{
		cards: make(map[uuid.UUID]*models.GameCard),
	}

	// Add some sample cards for development
	sampleCards := []*models.GameCard{}

	// Populate the mock storage with sample data
	for _, card := range sampleCards {
		storage.cards[card.ID] = card
	}

	return storage
}

// ListCards returns all cards in storage
func (m *MockStorage) ListCards(ctx context.Context) ([]*models.GameCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cards := make([]*models.GameCard, 0, len(m.cards))
	for _, card := range m.cards {
		// Create a copy to avoid modifying the original
		cardCopy := *card
		cards = append(cards, &cardCopy)
	}

	return cards, nil
}

// GetCard returns a specific card by ID
func (m *MockStorage) GetCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	card, exists := m.cards[id]
	if !exists {
		return nil, errors.New("card not found")
	}

	// Return a copy to avoid modifying the original
	cardCopy := *card
	return &cardCopy, nil
}

// CreateCard adds a new card to storage
func (m *MockStorage) CreateCard(ctx context.Context, card *models.GameCard) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a new ID if not provided
	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}

	// Check if card already exists
	if _, exists := m.cards[card.ID]; exists {
		return errors.New("card already exists")
	}

	// Store a copy to avoid external modifications
	cardCopy := *card
	m.cards[card.ID] = &cardCopy

	return nil
}

// UpdateCard updates an existing card in storage
func (m *MockStorage) UpdateCard(ctx context.Context, id uuid.UUID, card *models.GameCard) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if card exists
	if _, exists := m.cards[id]; !exists {
		return errors.New("card not found")
	}

	// Update the card ID to match the URL parameter
	card.ID = id

	// Store a copy to avoid external modifications
	cardCopy := *card
	m.cards[id] = &cardCopy

	return nil
}

// DeleteCard removes a card from storage
func (m *MockStorage) DeleteCard(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if card exists
	if _, exists := m.cards[id]; !exists {
		return errors.New("card not found")
	}

	delete(m.cards, id)
	return nil
}
