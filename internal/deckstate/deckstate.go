// Package deckstate models the runtime state of a deck during gameplay.
//
// Models in this package track the mutable aspects of a deck instance, such as card zones,
// card orientation, draw/discard behavior, and per-player hands. This state is
// created from an immutable deck template (see package deckdef), and evolves
// as players interact with the cards.
//
// The package makes no assumptions about specific game rules, but provides a
// flexible framework for modeling card positions, zone behavior, and visibility.
//
// This package is distinct from deckdef, which defines static deck and card templates.

package deckstate
