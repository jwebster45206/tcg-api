package deckstate

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jwebster45206/tcg-api/internal/deckdef"
)

const (
	FaceDown Facing = "face-down"
	FaceUp   Facing = "face-up"
	InHand   Facing = "in-hand"

	Upright  Orientation = "upright"
	Rotated  Orientation = "rotated"  // 90 degrees
	Inverted Orientation = "inverted" // 180 degrees

	ZoneTypeDraw      ZoneType = "draw"
	ZoneTypeDiscard   ZoneType = "discard"
	ZoneTypeHand      ZoneType = "hand"
	ZoneTypeTable     ZoneType = "table"
	ZoneTypeTemporary ZoneType = "temporary" // Transient zones

	// Reserved Zone Names
	ZoneNameDraw    = "draw"    // Draw pile, cards are face down
	ZoneNameDiscard = "discard" // Discard pile, cards are face up

	ZoneSizeUnlimited = 0 // Special size for unlimited zones
)

// ZoneItem represents something that can exist in a zone. In most cases, this will
// be a card. ZoneItem also allows for card groupings, such as a group of cards laid
// down together in a game.
type ZoneItem interface {
	Count() int
	GetCards() []CardInZone // Returns cards in this item
	GetID() uuid.UUID       // Get a card ID representing this item
}

// CardInZone represents a card with its state within a specific zone
type CardInZone struct {
	Card        deckdef.CardInterface `json:"card"`
	Facing      *Facing               `json:"facing,omitempty"` // If nil, inherits from zone
	Orientation *Orientation          `json:"orientation"`
	// Could add more game-specific properties like:
	// Tapped      bool          `json:"tapped,omitempty"`     // For games that use tapping
	// Counters    map[string]int `json:"counters,omitempty"`  // +1/+1 counters, etc.
}

func (c CardInZone) Count() int {
	return 1 // Each CardInZone represents a single card instance
}

func (c CardInZone) GetCards() []CardInZone {
	return []CardInZone{c} // Returns itself as the only card in this item
}

func (c CardInZone) GetID() uuid.UUID {
	return c.Card.GetID()
}

// assert that CardInZone implements ZoneItem interface
var _ ZoneItem = CardInZone{}

// GroupInZone represents a group of cards within a zone.
// This is mostly useful for games where cards can be laid down in groupings.
type GroupInZone struct {
	Name        string       `json:"name"`                  // Optional name for the group
	Cards       []CardInZone `json:"cards"`                 // Cards in this group
	Facing      *Facing      `json:"facing,omitempty"`      // If nil, cards inherit individually
	Orientation *Orientation `json:"orientation,omitempty"` // If nil, cards inherit individually
}

func (g GroupInZone) Count() int {
	return len(g.Cards)
}

func (g GroupInZone) GetCards() []CardInZone {
	resolvedCards := make([]CardInZone, len(g.Cards))
	for i, card := range g.Cards {
		resolvedCards[i] = card
		if card.Facing == nil && g.Facing != nil {
			resolvedCards[i].Facing = g.Facing
		}
		if card.Orientation == nil && g.Orientation != nil {
			resolvedCards[i].Orientation = g.Orientation
		}
	}
	return resolvedCards
}

func (g GroupInZone) GetID() uuid.UUID {
	if len(g.Cards) == 0 {
		return uuid.Nil
	}
	return g.Cards[0].GetID()
}

// assert that GroupInZone implements ZoneItem interface
var _ ZoneItem = GroupInZone{}

// Zone is a collection of cards and groups within a specific area of the game.
// Zones can represent different game states like draw piles, discard piles, hands, etc.
type Zone struct {
	Name          string     `json:"name"`            // "draw", "hand:1", "table:2", etc.
	Type          ZoneType   `json:"type"`            // Optional: "draw", "hand", "discard"
	DefaultFacing Facing     `json:"default_facing"`  // Default facing for cards in this zone
	Items         []ZoneItem `json:"items,omitempty"` // Ordered list of cards and groups
	// Owner TODO
}

// NewZone creates a new zone with default settings based on zone type
// For unknown/unlimited size, use ZoneSizeUnlimited (0).
func NewZone(name string, zoneType ZoneType, size int) Zone {
	var defaultFacing Facing

	// Set default facing based on zone type
	switch zoneType {
	case ZoneTypeDraw:
		defaultFacing = FaceDown // Draw pile cards are typically face down
	case ZoneTypeDiscard:
		defaultFacing = FaceUp // Discard pile cards are typically face up
	case ZoneTypeHand:
		defaultFacing = InHand // Hand cards are typically visible to owner
	case ZoneTypeTable:
		defaultFacing = FaceUp // Table cards are typically face up
	default:
		defaultFacing = FaceDown
	}

	return Zone{
		Name:          name,
		Type:          zoneType,
		DefaultFacing: defaultFacing,
		Items:         make([]ZoneItem, 0, size),
	}
}

func IsValidZoneType(zoneType ZoneType) bool {
	switch zoneType {
	case ZoneTypeDraw, ZoneTypeDiscard, ZoneTypeHand, ZoneTypeTable, ZoneTypeTemporary:
		return true
	default:
		return false
	}
}

// ZoneItemJSON represents a ZoneItem for JSON marshaling/unmarshaling
type ZoneItemJSON struct {
	Type string          `json:"type"` // "card" or "group"
	Data json.RawMessage `json:"data"`
}

// MarshalJSON custom marshaler for Zone
func (z Zone) MarshalJSON() ([]byte, error) {
	// Create a temporary struct for marshaling
	temp := struct {
		Name          string         `json:"name"`
		Type          ZoneType       `json:"type"`
		DefaultFacing Facing         `json:"default_facing"`
		Items         []ZoneItemJSON `json:"items,omitempty"`
	}{
		Name:          z.Name,
		Type:          z.Type,
		DefaultFacing: z.DefaultFacing,
	}

	// Only process items if the slice is not nil
	if z.Items != nil {
		// Convert ZoneItems to ZoneItemJSON
		itemsJSON := make([]ZoneItemJSON, len(z.Items))
		for i, item := range z.Items {
			var itemType string
			var itemData json.RawMessage
			var err error

			switch v := item.(type) {
			case CardInZone:
				itemType = "card"
				itemData, err = json.Marshal(v)
			case GroupInZone:
				itemType = "group"
				itemData, err = json.Marshal(v)
			default:
				return nil, fmt.Errorf("unknown ZoneItem type: %T", v)
			}

			if err != nil {
				return nil, err
			}

			itemsJSON[i] = ZoneItemJSON{
				Type: itemType,
				Data: itemData,
			}
		}
		temp.Items = itemsJSON
	}
	return json.Marshal(temp)
}

// UnmarshalJSON custom unmarshaler for Zone
func (z *Zone) UnmarshalJSON(data []byte) error {
	// Temporary struct for unmarshaling
	temp := struct {
		Name          string         `json:"name"`
		Type          ZoneType       `json:"type"`
		DefaultFacing Facing         `json:"default_facing"`
		Items         []ZoneItemJSON `json:"items"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Set basic fields
	z.Name = temp.Name
	z.Type = temp.Type
	z.DefaultFacing = temp.DefaultFacing

	// Convert ZoneItemJSON back to ZoneItems
	z.Items = make([]ZoneItem, len(temp.Items))
	for i, itemJSON := range temp.Items {
		switch itemJSON.Type {
		case "card":
			var card CardInZone
			if err := json.Unmarshal(itemJSON.Data, &card); err != nil {
				return fmt.Errorf("failed to unmarshal CardInZone: %w", err)
			}
			z.Items[i] = card
		case "group":
			var group GroupInZone
			if err := json.Unmarshal(itemJSON.Data, &group); err != nil {
				return fmt.Errorf("failed to unmarshal GroupInZone: %w", err)
			}
			z.Items[i] = group
		default:
			return fmt.Errorf("unknown ZoneItem type: %s", itemJSON.Type)
		}
	}

	return nil
}

// CardInZoneJSON represents a CardInZone for JSON marshaling/unmarshaling
type CardInZoneJSON struct {
	Card        json.RawMessage `json:"card"`
	Facing      *Facing         `json:"facing,omitempty"`
	Orientation *Orientation    `json:"orientation,omitempty"`
}

// MarshalJSON custom marshaler for CardInZone
func (c CardInZone) MarshalJSON() ([]byte, error) {
	// Marshal the card directly, then add card_type
	cardJSON, err := json.Marshal(c.Card)
	if err != nil {
		return nil, err
	}

	// Parse the marshaled card to add card_type
	var cardMap map[string]interface{}
	if err := json.Unmarshal(cardJSON, &cardMap); err != nil {
		return nil, err
	}

	// Add the card_type field
	cardMap["card_type"] = c.Card.GetCardType()

	// Re-marshal with card_type included
	cardWithTypeJSON, err := json.Marshal(cardMap)
	if err != nil {
		return nil, err
	}

	result := CardInZoneJSON{
		Card:        cardWithTypeJSON,
		Facing:      c.Facing,
		Orientation: c.Orientation,
	}

	return json.Marshal(result)
}

// UnmarshalJSON custom unmarshaler for CardInZone
func (c *CardInZone) UnmarshalJSON(data []byte) error {
	var temp CardInZoneJSON
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	c.Facing = temp.Facing
	c.Orientation = temp.Orientation

	// Parse the card type to determine which struct to unmarshal into
	var cardTypeInfo struct {
		CardType string `json:"card_type"`
	}
	if err := json.Unmarshal(temp.Card, &cardTypeInfo); err != nil {
		return err
	}

	// Unmarshal into the appropriate card type
	switch cardTypeInfo.CardType {
	case deckdef.TypePlayingCard:
		var card deckdef.PlayingCard
		if err := json.Unmarshal(temp.Card, &card); err != nil {
			return err
		}
		c.Card = &card
	case deckdef.TypeImageCard:
		var card deckdef.ImageCard
		if err := json.Unmarshal(temp.Card, &card); err != nil {
			return err
		}
		c.Card = &card
	case deckdef.TypeGameCard:
		var card deckdef.GameCard
		if err := json.Unmarshal(temp.Card, &card); err != nil {
			return err
		}
		c.Card = &card
	default:
		return fmt.Errorf("unknown card type: %s", cardTypeInfo.CardType)
	}

	return nil
}
