package state

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/pkg/deckstate"
)

// MockDeckStateStorage is a mock implementation of DeckStateStorage for testing
type MockDeckStateStorage struct {
	mu     sync.RWMutex
	states map[string]*deckstate.DeckState
	locks  map[string]*mockLock

	// For testing error scenarios
	saveError   error
	getError    error
	deleteError error
	lockError   error
}

// mockLock implements DeckStateLock for testing
type mockLock struct {
	gameID     string
	ownerID    string
	expiration time.Time
	storage    *MockDeckStateStorage
}

func (l *mockLock) Release() error {
	l.storage.mu.Lock()
	defer l.storage.mu.Unlock()
	delete(l.storage.locks, l.gameID)
	return nil
}

func (l *mockLock) IsExpired() bool {
	return time.Now().After(l.expiration)
}

func (l *mockLock) OwnerID() string {
	return l.ownerID
}

func (l *mockLock) Extend(duration time.Duration) error {
	l.expiration = time.Now().Add(duration)
	return nil
}

// NewMockDeckStateStorage creates a new mock storage
func NewMockDeckStateStorage() *MockDeckStateStorage {
	return &MockDeckStateStorage{
		states: make(map[string]*deckstate.DeckState),
		locks:  make(map[string]*mockLock),
	}
}

// Ensure that MockDeckStateStorage implements DeckStateStorage interface
var _ DeckStateStorage = (*MockDeckStateStorage)(nil)

func (m *MockDeckStateStorage) SaveDeckState(ctx context.Context, gameID string, state *deckstate.DeckState) error {
	if m.saveError != nil {
		return m.saveError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Make a copy of the state to avoid shared references
	stateCopy := *state
	m.states[gameID] = &stateCopy

	return nil
}

func (m *MockDeckStateStorage) GetDeckState(ctx context.Context, gameID string) (*deckstate.DeckState, error) {
	if m.getError != nil {
		return nil, m.getError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[gameID]
	if !exists {
		return nil, nil // No state found
	}

	// Return a copy to avoid shared references
	stateCopy := *state
	return &stateCopy, nil
}

func (m *MockDeckStateStorage) DeleteDeckState(ctx context.Context, gameID string) error {
	if m.deleteError != nil {
		return m.deleteError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.states[gameID]
	if !exists {
		return nil // Redis DEL returns success even if key doesn't exist
	}

	delete(m.states, gameID)
	return nil
}

// Test helper methods

// SetSaveError sets an error to be returned by SaveDeckState
func (m *MockDeckStateStorage) SetSaveError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saveError = err
}

// SetGetError sets an error to be returned by GetDeckState
func (m *MockDeckStateStorage) SetGetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getError = err
}

// SetDeleteError sets an error to be returned by DeleteDeckState
func (m *MockDeckStateStorage) SetDeleteError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteError = err
}

// ClearErrors clears all configured errors
func (m *MockDeckStateStorage) ClearErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.saveError = nil
	m.getError = nil
	m.deleteError = nil
}

// GetStoredStates returns all stored states (for testing)
func (m *MockDeckStateStorage) GetStoredStates() map[string]*deckstate.DeckState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy of the map
	states := make(map[string]*deckstate.DeckState)
	for k, v := range m.states {
		stateCopy := *v
		states[k] = &stateCopy
	}
	return states
}

// Clear removes all stored states
func (m *MockDeckStateStorage) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states = make(map[string]*deckstate.DeckState)
}

// Count returns the number of stored states
func (m *MockDeckStateStorage) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.states)
}

// HasState checks if a state exists for the given game ID
func (m *MockDeckStateStorage) HasState(gameID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.states[gameID]
	return exists
}

// LockDeckState attempts to acquire a lock on the deck state
func (m *MockDeckStateStorage) LockDeckState(ctx context.Context, gameID string, timeout time.Duration) (DeckStateLock, error) {
	if m.lockError != nil {
		return nil, m.lockError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already locked
	if lock, exists := m.locks[gameID]; exists && !lock.IsExpired() {
		return nil, &LockError{GameID: gameID, Message: "deck state is locked by another process"}
	}

	// Create new lock
	ownerID := uuid.New().String()
	lock := &mockLock{
		gameID:     gameID,
		ownerID:    ownerID,
		expiration: time.Now().Add(timeout),
		storage:    m,
	}
	m.locks[gameID] = lock
	return lock, nil
}

// IsLocked checks if a deck state is currently locked
func (m *MockDeckStateStorage) IsLocked(ctx context.Context, gameID string) (bool, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lock, exists := m.locks[gameID]
	if !exists || lock.IsExpired() {
		return false, "", nil
	}
	return true, lock.OwnerID(), nil
}
