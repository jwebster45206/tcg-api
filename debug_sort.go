package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
	"github.com/jwebster45206/tcg-api/internal/shuffle"
)

func main() {
	// Create test cards
	dragonID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	castleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
	phoenixID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")

	dragonCard := &deckdef.ImageCard{
		ID:            dragonID,
		Name:          "Dragon",
		FrontImageURL: "https://example.com/dragon_front.png",
		BackImageURL:  "https://example.com/dragon_back.png",
		Description:   "A mighty dragon",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	castleCard := &deckdef.ImageCard{
		ID:            castleID,
		Name:          "Castle",
		FrontImageURL: "https://example.com/castle_front.png",
		BackImageURL:  "https://example.com/castle_back.png",
		Description:   "A strong castle",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	phoenixCard := &deckdef.ImageCard{
		ID:            phoenixID,
		Name:          "Phoenix",
		FrontImageURL: "https://example.com/phoenix_front.png",
		BackImageURL:  "https://example.com/phoenix_back.png",
		Description:   "A majestic phoenix",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	deck := &deckdef.Deck{
		ID:        uuid.New(),
		Name:      "Test Deck",
		Type:      "standard",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Cards: &deckdef.CardCollection{
			TotalCount:  5, // 2 dragons + 1 castle + 2 phoenix
			UniqueCount: 3,
			Items: []deckdef.CardWithQuantity{
				{Card: dragonCard, Quantity: 2},  // Positions 0, 1
				{Card: castleCard, Quantity: 1},  // Position 2
				{Card: phoenixCard, Quantity: 2}, // Positions 3, 4
			},
		},
	}

	// Create zone items in wrong order: Phoenix, Castle, Dragon, Dragon, Phoenix
	zoneItems := []deckstate.ZoneItem{
		deckstate.CardInZone{Card: phoenixCard},
		deckstate.CardInZone{Card: castleCard},
		deckstate.CardInZone{Card: dragonCard},
		deckstate.CardInZone{Card: dragonCard},
		deckstate.CardInZone{Card: phoenixCard},
	}

	fmt.Println("Before sorting:")
	for i, item := range zoneItems {
		cardInZone := item.(deckstate.CardInZone)
		fmt.Printf("  %d: %s\n", i, cardInZone.Card.GetName())
	}

	// Debug: show expanded cards
	expanded := deck.ExpandCards()
	fmt.Println("\nExpanded deck:")
	for i, card := range expanded {
		fmt.Printf("  %d: %s (%s)\n", i, card.GetName(), card.GetID().String())
	}

	// Debug: show what positions are assigned
	expandedPositions := make(map[string][]int)
	for position, card := range expanded {
		cardID := card.GetID().String()
		expandedPositions[cardID] = append(expandedPositions[cardID], position)
	}

	fmt.Println("\nPosition assignments:")
	usedPositions := make(map[string]int)
	for i, item := range zoneItems {
		cardInZone := item.(deckstate.CardInZone)
		cardID := cardInZone.Card.GetID().String()

		positions := expandedPositions[cardID]
		usedIndex := usedPositions[cardID]
		position := 999999
		if usedIndex < len(positions) {
			position = positions[usedIndex]
			usedPositions[cardID] = usedIndex + 1
		}

		fmt.Printf("  Item %d (%s): position %d\n", i, cardInZone.Card.GetName(), position)
	}

	// Sort the zone items
	err := shuffle.DefinitionSort(zoneItems, deck)
	if err != nil {
		fmt.Printf("DefinitionSort failed: %v\n", err)
		return
	}

	fmt.Println("\nAfter sorting:")
	for i, item := range zoneItems {
		cardInZone := item.(deckstate.CardInZone)
		fmt.Printf("  %d: %s (%s)\n", i, cardInZone.Card.GetName(), cardInZone.Card.GetID().String())
	}

	// Expected order: Dragon, Dragon, Castle, Phoenix, Phoenix
	expectedOrder := []uuid.UUID{dragonID, dragonID, castleID, phoenixID, phoenixID}
	fmt.Println("\nExpected order:")
	for i, id := range expectedOrder {
		fmt.Printf("  %d: %s\n", i, id.String())
	}
}
