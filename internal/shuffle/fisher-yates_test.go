package shuffle

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

func TestFisherYatesShuffle(t *testing.T) {
	// Create an ordered list of ZoneItems for testing
	items := createOrderedZoneItems(10)
	
	// Make a copy to compare original order
	originalOrder := make([]deckstate.ZoneItem, len(items))
	copy(originalOrder, items)
	
	// Perform first shuffle
	err := FisherYatesShuffle(items)
	if err != nil {
		t.Fatalf("FisherYatesShuffle returned error: %v", err)
	}
	
	// Check that the order has changed (very high probability)
	sameOrder := true
	for i := range items {
		if !areZoneItemsEqual(items[i], originalOrder[i]) {
			sameOrder = false
			break
		}
	}
	
	if sameOrder {
		t.Error("Shuffle did not change the order (extremely unlikely but possible)")
	}
	
	// Make a copy of first shuffle result
	firstShuffle := make([]deckstate.ZoneItem, len(items))
	copy(firstShuffle, items)
	
	// Perform second shuffle
	err = FisherYatesShuffle(items)
	if err != nil {
		t.Fatalf("Second FisherYatesShuffle returned error: %v", err)
	}
	
	// Check that the two shuffles are different (very high probability)
	sameAsPrevious := true
	for i := range items {
		if !areZoneItemsEqual(items[i], firstShuffle[i]) {
			sameAsPrevious = false
			break
		}
	}
	
	if sameAsPrevious {
		t.Error("Two consecutive shuffles produced identical results (extremely unlikely but possible)")
	}
	
	// Verify all original items are still present (no items lost/duplicated)
	if len(items) != len(originalOrder) {
		t.Fatalf("Shuffle changed slice length: expected %d, got %d", len(originalOrder), len(items))
	}
	
	// Count items to ensure no duplicates/losses
	originalCounts := countZoneItems(originalOrder)
	shuffledCounts := countZoneItems(items)
	
	for cardID, originalCount := range originalCounts {
		shuffledCount := shuffledCounts[cardID]
		if shuffledCount != originalCount {
			t.Errorf("Card %s count changed: expected %d, got %d", cardID, originalCount, shuffledCount)
		}
	}
}

func TestFisherYatesShuffle_EdgeCases(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		var items []deckstate.ZoneItem
		err := FisherYatesShuffle(items)
		if err != nil {
			t.Errorf("FisherYatesShuffle with empty slice returned error: %v", err)
		}
	})
	
	t.Run("single item", func(t *testing.T) {
		items := createOrderedZoneItems(1)
		original := items[0]
		
		err := FisherYatesShuffle(items)
		if err != nil {
			t.Errorf("FisherYatesShuffle with single item returned error: %v", err)
		}
		
		if !areZoneItemsEqual(items[0], original) {
			t.Error("Single item was modified during shuffle")
		}
	})
	
	t.Run("two items", func(t *testing.T) {
		items := createOrderedZoneItems(2)
		
		err := FisherYatesShuffle(items)
		if err != nil {
			t.Errorf("FisherYatesShuffle with two items returned error: %v", err)
		}
		
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	})
}

// Helper function to create an ordered list of ZoneItems for testing
func createOrderedZoneItems(count int) []deckstate.ZoneItem {
	items := make([]deckstate.ZoneItem, count)
	
	for i := 0; i < count; i++ {
		card := &deckdef.ImageCard{
			ID:            uuid.New(),
			Name:          fmt.Sprintf("Card %d", i),
			Description:   fmt.Sprintf("Test card number %d", i),
			FrontImageURL: "https://example.com/front.jpg",
			BackImageURL:  "https://example.com/back.jpg",
		}
		
		items[i] = deckstate.CardInZone{
			Card:        card,
			Facing:      nil,
			Orientation: nil,
		}
	}
	
	return items
}

// Helper function to compare two ZoneItems for equality
func areZoneItemsEqual(a, b deckstate.ZoneItem) bool {
	// For CardInZone, compare the card IDs
	cardA, okA := a.(deckstate.CardInZone)
	cardB, okB := b.(deckstate.CardInZone)
	
	if okA && okB {
		return cardA.Card.GetID() == cardB.Card.GetID()
	}
	
	// For other types, this would need to be extended
	return false
}

// Helper function to count occurrences of each card ID
func countZoneItems(items []deckstate.ZoneItem) map[uuid.UUID]int {
	counts := make(map[uuid.UUID]int)
	
	for _, item := range items {
		if card, ok := item.(deckstate.CardInZone); ok {
			counts[card.Card.GetID()]++
		}
	}
	
	return counts
}
