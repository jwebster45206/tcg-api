package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/models"
	"github.com/jwebster45206/tcg-api/internal/query"
)

var (
	ErrNotFound = errors.New("not found")
)

// MockStorage implements Storage interface for testing and development
type MockStorage struct {
	mu           sync.RWMutex
	gameCards    map[uuid.UUID]*models.GameCard
	decks        map[uuid.UUID]*models.Deck
	imageCards   map[uuid.UUID]*models.ImageCard
	playingCards map[uuid.UUID]*models.PlayingCard
}

// NewMockStorage creates a new MockStorage instance with some sample data
func NewMockStorage() Storage {
	storage := &MockStorage{
		gameCards:    make(map[uuid.UUID]*models.GameCard),
		decks:        make(map[uuid.UUID]*models.Deck),
		imageCards:   make(map[uuid.UUID]*models.ImageCard),
		playingCards: make(map[uuid.UUID]*models.PlayingCard),
	}

	// Add some sample cards for development
	sampleCards := []*models.GameCard{}

	// Populate the mock storage with sample data
	for _, card := range sampleCards {
		storage.gameCards[card.ID] = card
	}
	return storage
}

// Ping implements health check for mock storage (always returns nil)
func (ms *MockStorage) Ping(ctx context.Context) error {
	return nil
}

func (m *MockStorage) ListGameCards(ctx context.Context, cardType string) ([]*models.GameCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch cardType {
	case "gamecard":
		cards := make([]*models.GameCard, 0, len(m.gameCards))
		for _, card := range m.gameCards {
			// Create a copy to avoid modifying the original
			cardCopy := *card
			cards = append(cards, &cardCopy)
		}
		return cards, nil
	default:
		return nil, errors.New("unsupported card type")
	}
}

func (m *MockStorage) GetGameCard(ctx context.Context, id uuid.UUID) (*models.GameCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	card, exists := m.gameCards[id]
	if !exists {
		return nil, ErrNotFound
	}
	// Return a copy to avoid modifying the original
	cardCopy := *card
	return &cardCopy, nil
}

// CreateGameCard adds a new card to storage
func (m *MockStorage) CreateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}
	if _, exists := m.gameCards[card.ID]; exists {
		return nil, errors.New("card already exists")
	}

	// Store a copy to avoid external modifications
	cardCopy := card
	m.gameCards[card.ID] = &cardCopy

	return &cardCopy, nil
}

func (m *MockStorage) UpdateGameCard(ctx context.Context, card models.GameCard) (*models.GameCard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.gameCards[card.ID]; !exists {
		return nil, ErrNotFound
	}

	// Store a copy to avoid external modifications
	cardCopy := card
	m.gameCards[card.ID] = &cardCopy
	return &cardCopy, nil
}

func (m *MockStorage) DeleteGameCard(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.gameCards[id]; !exists {
		return ErrNotFound
	}
	delete(m.gameCards, id)
	return nil
}

// Deck operations

func (m *MockStorage) ListDecks(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.Deck, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	decks := make([]*models.Deck, 0, len(m.decks))
	for _, deck := range m.decks {
		deckCopy := *deck
		decks = append(decks, &deckCopy)
	}

	// Note: Mock implementation doesn't actually apply filters/sorts/pagination
	// In a real test scenario, you would implement proper filtering logic
	// For now, just return all decks (useful for basic testing)

	return decks, nil
}

func (m *MockStorage) GetDeck(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	deck, exists := m.decks[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Return a copy to avoid modifying the original
	deckCopy := *deck
	return &deckCopy, nil
}

func (m *MockStorage) CreateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if deck.ID == uuid.Nil {
		deck.ID = uuid.New()
	}

	if _, exists := m.decks[deck.ID]; exists {
		return nil, errors.New("deck already exists")
	}

	// Store a copy to avoid external modifications
	deckCopy := deck
	m.decks[deck.ID] = &deckCopy
	return &deckCopy, nil
}

func (m *MockStorage) UpdateDeck(ctx context.Context, deck models.Deck) (*models.Deck, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.decks[deck.ID]; !exists {
		return nil, ErrNotFound
	}

	// Store a copy to avoid external modifications
	deckCopy := deck
	m.decks[deck.ID] = &deckCopy
	return &deckCopy, nil
}

func (m *MockStorage) DeleteDeck(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.decks[id]; !exists {
		return ErrNotFound
	}
	delete(m.decks, id)
	return nil
}

func (m *MockStorage) CreateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a new ID if not provided
	if imageCard.ID == uuid.Nil {
		imageCard.ID = uuid.New()
	}

	if _, exists := m.imageCards[imageCard.ID]; exists {
		return nil, errors.New("image card already exists")
	}

	// Store a copy to avoid external modifications
	imageCopy := imageCard
	m.imageCards[imageCard.ID] = &imageCopy
	return &imageCopy, nil
}

func (m *MockStorage) GetImageCard(ctx context.Context, id uuid.UUID) (*models.ImageCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	imageCard, exists := m.imageCards[id]
	if !exists {
		return nil, ErrNotFound
	}
	// Return a copy to avoid modifying the original
	imageCopy := *imageCard
	return &imageCopy, nil
}

func (m *MockStorage) UpdateImageCard(ctx context.Context, imageCard models.ImageCard) (*models.ImageCard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.imageCards[imageCard.ID]; !exists {
		return nil, ErrNotFound
	}

	// Store a copy to avoid external modifications
	imageCopy := imageCard
	m.imageCards[imageCard.ID] = &imageCopy
	return &imageCopy, nil
}

func (m *MockStorage) DeleteImageCard(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.imageCards[id]; !exists {
		return ErrNotFound
	}
	delete(m.imageCards, id)
	return nil
}

func (m *MockStorage) ListImageCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.ImageCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// filter and sort are not implemented in mock storage
	imageCards := make([]*models.ImageCard, 0, len(m.imageCards))
	for _, imageCard := range m.imageCards {
		// Create a copy to avoid modifying the original
		imageCopy := *imageCard
		imageCards = append(imageCards, &imageCopy)
	}
	return imageCards, nil
}

func (m *MockStorage) ListPlayingCards(ctx context.Context, filters []query.Filter, sorts []query.SortOption, pageSize int, pageNum int) ([]*models.PlayingCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	playingCards := make([]*models.PlayingCard, 0, len(m.playingCards))
	for _, playingCard := range m.playingCards {
		// Create a copy to avoid modifying the original
		cardCopy := *playingCard
		playingCards = append(playingCards, &cardCopy)
	}

	// TODO: Apply filters, sorts, and pagination like ListImageCards
	// For now, return all cards (mock implementation)
	return playingCards, nil
}

func (m *MockStorage) GetPlayingCard(ctx context.Context, id uuid.UUID) (*models.PlayingCard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	playingCard, exists := m.playingCards[id]
	if !exists {
		return nil, ErrNotFound
	}
	// Return a copy to avoid modifying the original
	cardCopy := *playingCard
	return &cardCopy, nil
}

func (m *MockStorage) CreatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if card.ID == uuid.Nil {
		card.ID = uuid.New()
	}
	if err := card.Validate(); err != nil {
		return nil, err
	}
	if _, exists := m.playingCards[card.ID]; exists {
		return nil, errors.New("playing card already exists")
	}

	// Store a copy to avoid external modifications
	cardCopy := card
	m.playingCards[card.ID] = &cardCopy
	return &cardCopy, nil
}

func (m *MockStorage) UpdatePlayingCard(ctx context.Context, card models.PlayingCard) (*models.PlayingCard, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := card.Validate(); err != nil {
		return nil, err
	}
	if _, exists := m.playingCards[card.ID]; !exists {
		return nil, ErrNotFound
	}

	// Store a copy to avoid external modifications
	cardCopy := card
	m.playingCards[card.ID] = &cardCopy
	return &cardCopy, nil
}

func (m *MockStorage) DeletePlayingCard(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.playingCards[id]; !exists {
		return ErrNotFound
	}
	delete(m.playingCards, id)
	return nil
}

func (m *MockStorage) ListDeckCards(ctx context.Context, deckID uuid.UUID) ([]*models.CardWithQuantity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if the deck exists
	if _, exists := m.decks[deckID]; !exists {
		return nil, ErrNotFound
	}

	// For mock storage, return some sample cards with quantities
	// In a real implementation, this would query the deck_cards relationship table
	cards := make([]*models.CardWithQuantity, 0)

	// Add some sample playing cards if any exist in storage
	count := 0
	for _, playingCard := range m.playingCards {
		if count >= 3 { // Limit to 3 cards for mock
			break
		}
		cardCopy := *playingCard
		cardWithQuantity := &models.CardWithQuantity{
			Card:     &cardCopy,
			Quantity: count + 1, // Sample quantities: 1, 2, 3
		}
		cards = append(cards, cardWithQuantity)
		count++
	}

	// Add some sample image cards if any exist in storage
	count = 0
	for _, imageCard := range m.imageCards {
		if count >= 2 { // Limit to 2 cards for mock
			break
		}
		cardCopy := *imageCard
		cardWithQuantity := &models.CardWithQuantity{
			Card:     &cardCopy,
			Quantity: count + 1, // Sample quantities: 1, 2
		}
		cards = append(cards, cardWithQuantity)
		count++
	}

	return cards, nil
}
