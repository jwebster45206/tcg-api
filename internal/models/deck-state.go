package models

type Facing string

const (
	FaceDown Facing = "face-down"
	FaceUp   Facing = "face-up"
	InHand   Facing = "in-hand"
)

type Orientation string

const (
	Upright  Orientation = "upright"
	Rotated  Orientation = "rotated"  // 90 degrees
	Inverted Orientation = "inverted" // 180 degrees
)

// CardInZone represents a card with its state within a specific zone
type CardInZone struct {
	Card        CardInterface `json:"card"`
	Facing      *Facing       `json:"facing,omitempty"` // If nil, inherits from zone
	Orientation *Orientation  `json:"orientation"`
	// Could add more game-specific properties like:
	// Tapped      bool          `json:"tapped,omitempty"`     // For games that use tapping
	// Counters    map[string]int `json:"counters,omitempty"`  // +1/+1 counters, etc.
}

type ZoneType string

const (
	ZoneTypeDraw    ZoneType = "draw"
	ZoneTypeDiscard ZoneType = "discard"
	ZoneTypeHand    ZoneType = "hand"
	ZoneTypeTable   ZoneType = "table"
)

type Zone struct {
	Name          string       `json:"name"`           // "draw", "hand:1", "table:2", etc.
	Type          ZoneType     `json:"type"`           // Optional: "draw", "hand", "discard"
	DefaultFacing Facing       `json:"default_facing"` // Default facing for cards in this zone
	Cards         []CardInZone `json:"cards"`          // Ordered list of card states
	// Owner TODO
}

type DeckState struct {
	ID          string          `json:"id"`           // Unique ID for the state instance
	Deck        Deck            `json:"deck"`         // Source deck template
	PlayerCount int             `json:"player_count"` // For initializing zone layout
	Zones       map[string]Zone `json:"zones"`        // e.g. "draw", "discard", "hand:1"
}
