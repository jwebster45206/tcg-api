package shuffle

import (
	"math/rand"
	"time"

	"github.com/jwebster45206/tcg-api/pkg/deckstate"
)

// FisherYatesShuffle performs an in-place Fisher-Yates shuffle on a slice of ZoneItems.
func FisherYatesShuffle(items []deckstate.ZoneItem) error {
	if len(items) <= 1 {
		return nil // Nothing to shuffle
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Start from the last element and work backwards.
	// This gives each element exactly one chance to be swapped.
	for i := len(items) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		items[i], items[j] = items[j], items[i]
	}
	return nil
}
