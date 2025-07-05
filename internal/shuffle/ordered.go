package shuffle

import (
	"github.com/jwebster45206/tcg-api/internal/deckdef"
	"github.com/jwebster45206/tcg-api/internal/deckstate"
)

// Definition sort attempts to sort zoneItems in the order of cards
// native to the deck definition.
// If the zoneItem is a group, the sorting logic uses the first card in the group.
//
// We store two floating indexes, one for zoneItems and one for the deck definition.
// The sort iterates through deckdef cards, but exits early if all zoneItems are sorted.
// For each definition card, it finds the first zoneItem that matches.
// (If a zoneitem is a group, it recursively sorts the group; then it uses the group's first card.)
// If it finds a match, it moves that zoneItem to the front of the sorted list.
// If it doesn't find a match, it continues to the next definition card.
func DefinitionSort(zoneItems []deckstate.ZoneItem, deck *deckdef.Deck) error {

}
