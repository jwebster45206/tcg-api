package storage

import (
	"context"
	"errors"
	"fmt"
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
	deckCards    map[uuid.UUID][]*models.CardWithQuantity // deckID -> cards
}

// NewMockStorage creates a new MockStorage instance with some sample data
func NewMockStorage() Storage {
	storage := &MockStorage{
		gameCards:    make(map[uuid.UUID]*models.GameCard),
		decks:        make(map[uuid.UUID]*models.Deck),
		imageCards:   make(map[uuid.UUID]*models.ImageCard),
		playingCards: make(map[uuid.UUID]*models.PlayingCard),
		deckCards:    make(map[uuid.UUID][]*models.CardWithQuantity),
	}

	// Add some sample playing cards for testing
	samplePlayingCards := []*models.PlayingCard{
		{
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Name:          "Ace of Spades",
			Description:   "The ace of spades playing card",
			FrontImageURL: "https://example.com/ace-spades.jpg",
			BackImageURL:  "https://example.com/back.jpg",
			Suit:          "spades",
			Ranking:       1,
		},
		{
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:          "King of Hearts",
			Description:   "The king of hearts playing card",
			FrontImageURL: "https://example.com/king-hearts.jpg",
			BackImageURL:  "https://example.com/back.jpg",
			Suit:          "hearts",
			Ranking:       13,
		},
		{
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:          "Queen of Diamonds",
			Description:   "The queen of diamonds playing card",
			FrontImageURL: "https://example.com/queen-diamonds.jpg",
			BackImageURL:  "https://example.com/back.jpg",
			Suit:          "diamonds",
			Ranking:       12,
		},
	}

	// Add some sample image cards
	sampleImageCards := []*models.ImageCard{
		{
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Name:          "Sample Image Card",
			Description:   "A sample image card for testing",
			FrontImageURL: "https://example.com/sample-front.jpg",
			BackImageURL:  "https://example.com/sample-back.jpg",
		},
	}

	// Add some sample game cards
	sampleGameCards := []*models.GameCard{
		{
			ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Name:          "Sample Game Card",
			FrontImageURL: "https://example.com/game-front.jpg",
			BackImageURL:  "https://example.com/game-back.jpg",
		},
	}

	// Populate the mock storage with sample data
	for _, card := range samplePlayingCards {
		storage.playingCards[card.ID] = card
	}
	for _, card := range sampleImageCards {
		storage.imageCards[card.ID] = card
	}
	for _, card := range sampleGameCards {
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

	// Return stored cards for this deck if they exist
	if cards, exists := m.deckCards[deckID]; exists {
		// Return a copy to prevent external modification
		result := make([]*models.CardWithQuantity, len(cards))
		for i, card := range cards {
			cardCopy := *card
			result[i] = &cardCopy
		}
		return result, nil
	}

	// If no cards are explicitly set, return empty slice
	return []*models.CardWithQuantity{}, nil
}

func (m *MockStorage) SetDeckCards(ctx context.Context, deckID uuid.UUID, cards []models.CardInputWithQuantity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.decks[deckID]; !exists {
		return ErrNotFound
	}

	// Validate that all cards exist and convert to CardWithQuantity
	cardPtrs := make([]*models.CardWithQuantity, 0, len(cards))

	for _, cardInput := range cards {
		var foundCard models.CardInterface

		for _, pc := range m.playingCards {
			if pc.ID == cardInput.Card.ID {
				foundCard = pc
				break
			}
		}

		if foundCard == nil {
			for _, ic := range m.imageCards {
				if ic.ID == cardInput.Card.ID {
					foundCard = ic
					break
				}
			}
		}

		if foundCard == nil {
			for _, gc := range m.gameCards {
				if gc.ID == cardInput.Card.ID {
					foundCard = gc
					break
				}
			}
		}

		if foundCard == nil {
			return fmt.Errorf("card not found: %s", cardInput.Card.ID)
		}

		cardWithQuantity := &models.CardWithQuantity{
			Card:     foundCard,
			Quantity: cardInput.Quantity,
		}
		cardPtrs = append(cardPtrs, cardWithQuantity)
	}

	m.deckCards[deckID] = cardPtrs
	return nil
}

func (m *MockStorage) GetCardsByIDs(ctx context.Context, cardIDs []uuid.UUID) ([]models.CardInterface, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var cards []models.CardInterface

	for _, cardID := range cardIDs {
		var found bool

		if card, exists := m.playingCards[cardID]; exists {
			cardCopy := *card
			cards = append(cards, &cardCopy)
			found = true
		}

		if !found {
			if card, exists := m.imageCards[cardID]; exists {
				cardCopy := *card
				cards = append(cards, &cardCopy)
				found = true
			}
		}

		if !found {
			if card, exists := m.gameCards[cardID]; exists {
				cardCopy := *card
				cards = append(cards, &cardCopy)
				found = true
			}
		}

		if !found {
			return nil, ErrNotFound
		}
	}

	return cards, nil
}
