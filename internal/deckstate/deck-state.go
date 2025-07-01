package deckstate

import "github.com/jwebster45206/tcg-api/internal/deckdef"

type Facing string
type Orientation string
type ZoneType string

// DeckState is a runtime state of a deck during gameplay.
// It includes the deck template, player count, and zones where cards are located.
type DeckState struct {
	ID          string          `json:"id"`           // Unique ID for the state instance
	Deck        deckdef.Deck    `json:"deck"`         // Source deck template
	PlayerCount int             `json:"player_count"` // For initializing zone layout
	Zones       map[string]Zone `json:"zones"`        // e.g. "draw", "discard", "hand:1"
}

// NewDeckState initializes a new DeckState from a Deck template.
// It sets up the initial zones based on the deck's cards.
func NewDeckState(deck deckdef.Deck, playerCount int) *DeckState {
	zones := make(map[string]Zone)

	if deck.Cards != nil {
		var drawItems []ZoneItem
		for _, cardWithQty := range deck.Cards.Items {
			for i := 0; i < cardWithQty.Quantity; i++ {
				cardInZone := CardInZone{
					Card:        cardWithQty.Card,
					Facing:      nil, // zone default
					Orientation: nil, // zone default
				}
				drawItems = append(drawItems, cardInZone)
			}
		}

		drawZone := NewZone(ZoneNameDraw, ZoneTypeDraw)
		drawZone.Items = drawItems
		zones[ZoneNameDraw] = drawZone
		zones[ZoneNameDiscard] = NewZone(ZoneNameDiscard, ZoneTypeDiscard)
	}

	return &DeckState{
		ID:          deck.ID.String(),
		Deck:        deck,
		PlayerCount: playerCount,
		Zones:       zones,
	}
}
