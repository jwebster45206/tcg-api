// Package deckdef contains definitions for card decks and card types.
// Game-specific card types may also live in this package, as long as they are
// purely definitional.
//
// Models in the package are immutable templates that define deck composition
// and card properties, but do not track gameplay state or card positions.
// These templates can be used to instantiate one or more independent DeckState instances.
//
// This package is distinct from deckstate, which models mutable gameplay
// state and card positions.

package deckdef
