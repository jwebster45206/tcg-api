package deckstate

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/pkg/deckdef"
)

type Facing string
type Orientation string
type ZoneType string

const (
	DefaultHandSize = 7
)

// DeckState is a runtime state of a deck during gameplay.
// It includes the deck template, player count, and zones where cards are located.
type DeckState struct {
	ID          string       `json:"id"`           // Unique ID for the state instance
	Deck        deckdef.Deck `json:"deck"`         // Source deck template
	PlayerCount int          `json:"player_count"` // For initializing zone layout

	// Zones uses pointers to allow direct mutation of zone fields.
	// This avoids the copy-modify-reassign pattern required with struct values in maps.
	Zones map[string]*Zone `json:"zones"` // e.g. "draw", "discard", "hand:1"
}

// NewDeckState initializes a new DeckState from a Deck template.
// It sets up the initial zones based on the deck's cards.
func NewDeckState(deck deckdef.Deck, playerCount int) *DeckState {
	zones := make(map[string]*Zone)

	if deck.Cards != nil {
		var drawItems []ZoneItem
		expandedCards := deck.ExpandCards()
		for _, card := range expandedCards {
			cardInZone := CardInZone{
				Card:        card,
				Facing:      nil, // zone default
				Orientation: nil, // zone default
			}
			drawItems = append(drawItems, cardInZone)
		}

		drawZone := NewZone(ZoneNameDraw, ZoneTypeDraw, deck.Cards.TotalCount)
		drawZone.Items = drawItems
		zones[ZoneNameDraw] = &drawZone

		discardZone := NewZone(ZoneNameDiscard, ZoneTypeDiscard, ZoneSizeUnlimited)
		zones[ZoneNameDiscard] = &discardZone

		if playerCount > 0 {
			for i := 1; i <= playerCount; i++ {
				handZone := NewZone(fmt.Sprintf("player:%d", i), ZoneTypeHand, DefaultHandSize)
				zones[handZone.Name] = &handZone
			}
		}
	}

	return &DeckState{
		ID:          uuid.New().String(),
		Deck:        deck,
		PlayerCount: playerCount,
		Zones:       zones,
	}
}
