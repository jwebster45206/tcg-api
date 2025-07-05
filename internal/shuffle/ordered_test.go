package shuffle

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

func TestDefinitionSort(t *testing.T) {
	// Test case 1: Standard deck of 54 cards, out of order
	t.Run("Standard54CardDeck", func(t *testing.T) {
		// Create a standard 54-card deck (52 cards + 2 jokers)
		deck := createStandard54CardDeck()

		// Create zone items in reverse order (worst case scenario)
		expandedCards := deck.ExpandCards()
		var zoneItems []deckstate.ZoneItem
		for i := len(expandedCards) - 1; i >= 0; i-- {
			zoneItems = append(zoneItems, deckstate.CardInZone{
				Card: expandedCards[i],
			})
		}

		// Verify they start out of order
		if zoneItems[0].(deckstate.CardInZone).Card.GetID() == expandedCards[0].GetID() {
			t.Fatal("Cards should start out of order for this test")
		}

		// Sort the zone items
		err := DefinitionSort(zoneItems, deck)
		if err != nil {
			t.Fatalf("DefinitionSort failed: %v", err)
		}

		// Verify they are now in correct order
		if len(zoneItems) != len(expandedCards) {
			t.Fatalf("Expected %d cards, got %d", len(expandedCards), len(zoneItems))
		}

		for i, item := range zoneItems {
			cardInZone := item.(deckstate.CardInZone)
			expectedCard := expandedCards[i]
			if cardInZone.Card.GetID() != expectedCard.GetID() {
				t.Errorf("Position %d: expected card %s, got %s",
					i, expectedCard.GetID(), cardInZone.Card.GetID())
			}
		}
	})

	// Test case 2: Deck made up of our 3 test image cards
	t.Run("ThreeImageCardDeck", func(t *testing.T) {
		// Create deck with our test cards
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

		// Sort the zone items
		err := DefinitionSort(zoneItems, deck)
		if err != nil {
			t.Fatalf("DefinitionSort failed: %v", err)
		}

		// Verify correct order: Dragon, Dragon, Castle, Phoenix, Phoenix
		expectedOrder := []uuid.UUID{dragonID, dragonID, castleID, phoenixID, phoenixID}

		if len(zoneItems) != len(expectedOrder) {
			t.Fatalf("Expected %d cards, got %d", len(expectedOrder), len(zoneItems))
		}

		for i, item := range zoneItems {
			cardInZone := item.(deckstate.CardInZone)
			expectedID := expectedOrder[i]
			if cardInZone.Card.GetID() != expectedID {
				t.Errorf("Position %d: expected card %s, got %s",
					i, expectedID, cardInZone.Card.GetID())
			}
		}
	})

	// Test case 3: Stack of cards partially matched with deck definition
	t.Run("PartialMatch", func(t *testing.T) {
		// Create a deck with 2 dragons and 1 castle
		dragonID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
		castleID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
		unknownID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440099")

		dragonCard := &deckdef.ImageCard{
			ID:            dragonID,
			Name:          "Dragon",
			FrontImageURL: "https://example.com/dragon_front.png",
			BackImageURL:  "https://example.com/dragon_back.png",
		}

		castleCard := &deckdef.ImageCard{
			ID:            castleID,
			Name:          "Castle",
			FrontImageURL: "https://example.com/castle_front.png",
			BackImageURL:  "https://example.com/castle_back.png",
		}

		unknownCard := &deckdef.ImageCard{
			ID:            unknownID,
			Name:          "Unknown Card",
			FrontImageURL: "https://example.com/unknown_front.png",
			BackImageURL:  "https://example.com/unknown_back.png",
		}

		deck := &deckdef.Deck{
			ID:        uuid.New(),
			Name:      "Partial Deck",
			Type:      "standard",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Cards: &deckdef.CardCollection{
				TotalCount:  3,
				UniqueCount: 2,
				Items: []deckdef.CardWithQuantity{
					{Card: dragonCard, Quantity: 2}, // Positions 0, 1
					{Card: castleCard, Quantity: 1}, // Position 2
				},
			},
		}

		// Create zone items: some match deck, some don't
		// Input: Unknown, Castle, Dragon, Unknown, Dragon, Unknown
		zoneItems := []deckstate.ZoneItem{
			deckstate.CardInZone{Card: unknownCard}, // Not in deck
			deckstate.CardInZone{Card: castleCard},  // Position 2 in deck
			deckstate.CardInZone{Card: dragonCard},  // Position 0 in deck
			deckstate.CardInZone{Card: unknownCard}, // Not in deck
			deckstate.CardInZone{Card: dragonCard},  // Position 1 in deck
			deckstate.CardInZone{Card: unknownCard}, // Not in deck
		}

		// Sort the zone items
		err := DefinitionSort(zoneItems, deck)
		if err != nil {
			t.Fatalf("DefinitionSort failed: %v", err)
		}

		// Verify:
		// - First 3 positions should be: Dragon, Dragon, Castle (in deck order)
		// - Last 3 positions should be: Unknown cards (moved to end, stable order)
		if len(zoneItems) != 6 {
			t.Fatalf("Expected 6 cards, got %d", len(zoneItems))
		}

		// Check first 3 positions (deck-defined cards)
		expectedDeckOrder := []uuid.UUID{dragonID, dragonID, castleID}
		for i := 0; i < 3; i++ {
			cardInZone := zoneItems[i].(deckstate.CardInZone)
			expectedID := expectedDeckOrder[i]
			if cardInZone.Card.GetID() != expectedID {
				t.Errorf("Position %d: expected deck card %s, got %s",
					i, expectedID, cardInZone.Card.GetID())
			}
		}

		// Check last 3 positions (unknown cards moved to end)
		for i := 3; i < 6; i++ {
			cardInZone := zoneItems[i].(deckstate.CardInZone)
			if cardInZone.Card.GetID() != unknownID {
				t.Errorf("Position %d: expected unknown card %s, got %s",
					i, unknownID, cardInZone.Card.GetID())
			}
		}
	})
}

// Helper function to convert rank string to integer ranking
func getRankingFromString(rank string) int {
	switch rank {
	case "ace":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "jack":
		return 11
	case "queen":
		return 12
	case "king":
		return 13
	default:
		return 0
	}
}

// Helper function to create a standard 54-card deck for testing
// This matches the database seed structure from full-deck-seeds.sql
func createStandard54CardDeck() *deckdef.Deck {
	var cards []deckdef.CardWithQuantity

	// Standard 52-card deck (4 suits × 13 ranks) - matches the database seed order
	suits := []string{"spades", "hearts", "diamonds", "clubs"}
	ranks := []string{"ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "jack", "queen", "king"}

	for _, suit := range suits {
		for _, rank := range ranks {
			card := &deckdef.PlayingCard{
				ID:            uuid.New(),
				Name:          rank + " of " + suit,
				Suit:          suit,
				Ranking:       getRankingFromString(rank),
				FrontImageURL: "https://example.com/" + rank + "_" + suit + ".png",
				BackImageURL:  "https://example.com/back.png",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			cards = append(cards, deckdef.CardWithQuantity{
				Card:     card,
				Quantity: 1,
			})
		}
	}

	// Add 2 jokers as ImageCards (matching the database seed)
	redJoker := &deckdef.ImageCard{
		ID:            uuid.New(),
		Name:          "Red Joker",
		FrontImageURL: "https://example.com/red_joker.png",
		BackImageURL:  "https://example.com/back.png",
		Description:   "Red Wild Card",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	blackJoker := &deckdef.ImageCard{
		ID:            uuid.New(),
		Name:          "Black Joker",
		FrontImageURL: "https://example.com/black_joker.png",
		BackImageURL:  "https://example.com/back.png",
		Description:   "Black Wild Card",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	cards = append(cards,
		deckdef.CardWithQuantity{Card: redJoker, Quantity: 1},
		deckdef.CardWithQuantity{Card: blackJoker, Quantity: 1})

	return &deckdef.Deck{
		ID:        uuid.New(),
		Name:      "Standard 54-Card Deck",
		Type:      "playing-card",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Cards: &deckdef.CardCollection{
			TotalCount:  54,
			UniqueCount: 54,
			Items:       cards,
		},
	}
}
