package models

const (
	FaceDown Facing = "face-down"
	FaceUp   Facing = "face-up"
	InHand   Facing = "in-hand"

	Upright  Orientation = "upright"
	Rotated  Orientation = "rotated"  // 90 degrees
	Inverted Orientation = "inverted" // 180 degrees

	ZoneTypeDraw    ZoneType = "draw"
	ZoneTypeDiscard ZoneType = "discard"
	ZoneTypeHand    ZoneType = "hand"
	ZoneTypeTable   ZoneType = "table"

	// Reserved Zone Names
	ZoneNameDraw    = "draw"    // Draw pile, cards are face down
	ZoneNameDiscard = "discard" // Discard pile, cards are face up
)

// ZoneItem represents something that can exist in a zone. In most cases, this will
// be a card. ZoneItem also allows for card groupings, such as a group of cards laid
// down together in a game.
type ZoneItem interface {
	Count() int
	GetCards() []CardInZone // Returns cards in this item
}

// CardInZone represents a card with its state within a specific zone
type CardInZone struct {
	Card        CardInterface `json:"card"`
	Facing      *Facing       `json:"facing,omitempty"` // If nil, inherits from zone
	Orientation *Orientation  `json:"orientation"`
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

type Zone struct {
	Name          string     `json:"name"`           // "draw", "hand:1", "table:2", etc.
	Type          ZoneType   `json:"type"`           // Optional: "draw", "hand", "discard"
	DefaultFacing Facing     `json:"default_facing"` // Default facing for cards in this zone
	Items         []ZoneItem `json:"items"`          // Ordered list of cards and groups
	// Owner TODO
}

// NewZone creates a new zone with default settings based on zone type
func NewZone(name string, zoneType ZoneType) Zone {
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
		Items:         make([]ZoneItem, 0),
	}
}
