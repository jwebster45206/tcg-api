package shuffle

import (
	"sort"

	"github.com/jwebster45206/tcg-api/pkg/deckdef"
	"github.com/jwebster45206/tcg-api/pkg/deckstate"
)

// DefinitionSort sorts zoneItems in the order of cards native to the deck definition.
// Groups are positioned by their first card.
// Complexity is O(n log n) for the sort.
func DefinitionSort(zoneItems []deckstate.ZoneItem, deck *deckdef.Deck) error {
	expandedCards := deck.ExpandCards()
	expandedPositions := make(map[string][]int) // cardID -> [position1, position2, ...]
	for position, card := range expandedCards {
		cardID := card.GetID().String()
		expandedPositions[cardID] = append(expandedPositions[cardID], position)
	}

	// Create a struct to hold the item and its position
	type itemWithPosition struct {
		item     deckstate.ZoneItem
		position int
	}

	// Create slice of items with their positions
	itemsWithPositions := make([]itemWithPosition, len(zoneItems))
	usedPositions := make(map[string]int) // cardID -> nextAvailableIndex

	for i, item := range zoneItems {
		position := getNextDefinitionPosition(item, expandedPositions, usedPositions)
		itemsWithPositions[i] = itemWithPosition{item: item, position: position}
	}

	// Sort by position
	sort.SliceStable(itemsWithPositions, func(i, j int) bool {
		return itemsWithPositions[i].position < itemsWithPositions[j].position
	})

	// Copy back to original slice
	for i, itemWithPos := range itemsWithPositions {
		zoneItems[i] = itemWithPos.item
	}

	return nil
}

// getNextDefinitionPosition returns the next available position for an item in the deck definition.
// Returns a large number for items not in the definition (moves them to end).
// Updates usedPositions to track which card instances have been assigned.
func getNextDefinitionPosition(item deckstate.ZoneItem, expandedPositions map[string][]int, usedPositions map[string]int) int {
	const notFoundPosition = 999999
	var cardID string

	switch v := item.(type) {
	case deckstate.CardInZone:
		cardID = v.Card.GetID().String()
	case deckstate.GroupInZone:
		if len(v.Cards) > 0 {
			cardID = v.Cards[0].Card.GetID().String()
		} else {
			return notFoundPosition
		}
	default:
		return notFoundPosition
	}

	// Check if this card exists in the deck definition
	positions, exists := expandedPositions[cardID]
	if !exists {
		return notFoundPosition
	}
	// Get the next available position for this card
	usedIndex := usedPositions[cardID]
	if usedIndex >= len(positions) {
		return notFoundPosition
	}
	// Mark this position as used and return it
	position := positions[usedIndex]
	usedPositions[cardID] = usedIndex + 1
	return position
}
