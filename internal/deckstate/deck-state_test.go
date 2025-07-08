package deckstate

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
)

func TestNewDeckState(t *testing.T) {
	// Create sample ImageCards matching the ones from seeds.sql
	dragonCard := &deckdef.ImageCard{
		ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Name:          "Dragon Artwork",
		Description:   "Beautiful dragon illustration card",
		FrontImageURL: "https://example.com/images/dragon-front.jpg",
		BackImageURL:  "https://example.com/images/card-back.jpg",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	castleCard := &deckdef.ImageCard{
		ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
		Name:          "Castle Scene",
		Description:   "Medieval castle landscape artwork",
		FrontImageURL: "https://example.com/images/castle-front.jpg",
		BackImageURL:  "https://example.com/images/card-back.jpg",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	phoenixCard := &deckdef.ImageCard{
		ID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
		Name:          "Phoenix Rising",
		Description:   "Majestic phoenix in flames",
		FrontImageURL: "https://example.com/images/phoenix-front.jpg",
		BackImageURL:  "https://example.com/images/card-back.jpg",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create a deck with these cards
	deckID := uuid.New()
	deck := deckdef.Deck{
		ID:        deckID,
		Name:      "Test Deck",
		Type:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Cards: &deckdef.CardCollection{
			TotalCount:  5, // 2 dragons + 1 castle + 2 phoenix = 5 total
			UniqueCount: 3, // 3 different cards
			Items: []deckdef.CardWithQuantity{
				{Card: dragonCard, Quantity: 2},  // 2 dragon cards
				{Card: castleCard, Quantity: 1},  // 1 castle card
				{Card: phoenixCard, Quantity: 2}, // 2 phoenix cards
			},
		},
	}

	// Test NewDeckState
	deckState := NewDeckState(deck, 2)

	// Verify basic properties
	if deckState.ID == uuid.Nil.String() {
		t.Errorf("Expected not nil, got %s", deckState.ID)
	}

	if deckState.PlayerCount != 2 {
		t.Errorf("Expected PlayerCount 2, got %d", deckState.PlayerCount)
	}

	if deckState.Deck.ID != deckID {
		t.Errorf("Expected Deck.ID %s, got %s", deckID, deckState.Deck.ID)
	}

	// Verify draw zone was created
	drawZone, exists := deckState.Zones["draw"]
	if !exists {
		t.Fatal("Expected draw zone to exist")
	}

	if drawZone.Name != "draw" {
		t.Errorf("Expected zone name 'draw', got '%s'", drawZone.Name)
	}

	if drawZone.Type != ZoneTypeDraw {
		t.Errorf("Expected zone type %s, got %s", ZoneTypeDraw, drawZone.Type)
	}

	if drawZone.DefaultFacing != FaceDown {
		t.Errorf("Expected default facing %s, got %s", FaceDown, drawZone.DefaultFacing)
	}

	// Verify player 1 hand was created
	_, exists = deckState.Zones["player:1"]
	if !exists {
		t.Fatal("Expected player:1 hand zone to exist")
	}

	// Verify player 2 hand was created
	_, exists = deckState.Zones["player:2"]
	if !exists {
		t.Fatal("Expected player:2 hand zone to exist")
	}

	// Verify all cards were expanded properly
	if len(drawZone.Items) != 5 {
		t.Errorf("Expected 5 cards in draw zone, got %d", len(drawZone.Items))
	}

	// Verify each item is a CardInZone
	for i, item := range drawZone.Items {
		cardInZone, ok := item.(CardInZone)
		if !ok {
			t.Errorf("Expected item %d to be CardInZone, got %T", i, item)
			continue
		}

		if cardInZone.Facing != nil {
			t.Errorf("Expected card %d facing to be nil (inherit from zone), got %v", i, cardInZone.Facing)
		}

		if cardInZone.Orientation != nil {
			t.Errorf("Expected card %d orientation to be nil (inherit from zone), got %v", i, cardInZone.Orientation)
		}

		// Verify it's one of our expected cards
		cardID := cardInZone.Card.GetID()
		if cardID != dragonCard.ID && cardID != castleCard.ID && cardID != phoenixCard.ID {
			t.Errorf("Unexpected card ID %s in draw zone", cardID)
		}
	}

	// Count occurrences of each card type
	dragonCount := 0
	castleCount := 0
	phoenixCount := 0

	for _, item := range drawZone.Items {
		cardInZone := item.(CardInZone)
		cardID := cardInZone.Card.GetID()
		switch cardID {
		case dragonCard.ID:
			dragonCount++
		case castleCard.ID:
			castleCount++
		case phoenixCard.ID:
			phoenixCount++
		}
	}

	if dragonCount != 2 {
		t.Errorf("Expected 2 dragon cards, got %d", dragonCount)
	}
	if castleCount != 1 {
		t.Errorf("Expected 1 castle card, got %d", castleCount)
	}
	if phoenixCount != 2 {
		t.Errorf("Expected 2 phoenix cards, got %d", phoenixCount)
	}

	// Test GetCards() method on ZoneItems
	for i, item := range drawZone.Items {
		// Each item should be a CardInZone, so GetCards() should return exactly one card
		cards := item.GetCards()
		if len(cards) != 1 {
			t.Errorf("Expected GetCards() to return 1 card for item %d, got %d", i, len(cards))
			continue
		}

		// The returned card should be the same as the original card
		returnedCard := cards[0]
		originalCard := item.(CardInZone)

		if returnedCard.Card.GetID() != originalCard.Card.GetID() {
			t.Errorf("GetCards() returned different card ID for item %d: expected %s, got %s",
				i, originalCard.Card.GetID(), returnedCard.Card.GetID())
		}

		// Check that facing/orientation inheritance works (should be nil since zone inheritance isn't applied yet)
		if returnedCard.Facing != nil {
			t.Errorf("Expected GetCards() card %d facing to be nil (not resolved), got %v", i, returnedCard.Facing)
		}
		if returnedCard.Orientation != nil {
			t.Errorf("Expected GetCards() card %d orientation to be nil (not resolved), got %v", i, returnedCard.Orientation)
		}
	}

	// Test Count() method on ZoneItems
	totalCount := 0
	for _, item := range drawZone.Items {
		count := item.Count()
		if count != 1 {
			t.Errorf("Expected Count() to return 1 for CardInZone, got %d", count)
		}
		totalCount += count
	}

	if totalCount != 5 {
		t.Errorf("Expected total count of all items to be 5, got %d", totalCount)
	}
}

func TestNewDeckState_EmptyDeck(t *testing.T) {
	// Test with a deck that has no cards
	deckID := uuid.New()
	deck := deckdef.Deck{
		ID:        deckID,
		Name:      "Empty Deck",
		Type:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Cards:     nil, // No cards
	}

	deckState := NewDeckState(deck, 1)

	// Should still create the deck state but with no zones
	if len(deckState.Zones) != 0 {
		t.Errorf("Expected 0 zones for empty deck, got %d", len(deckState.Zones))
	}
}
