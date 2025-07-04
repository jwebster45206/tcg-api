# tcg-api
A lightweight REST API for simulating card decks, built with Go. The API features a modular card interface design that separates game-specific mechanics from general deck mechanics, making it extensible for different types of card games.

## Technical Stack
- **Language**: Go
- **Storage**: MySQL (Redis TODO)

## Features

### Card Management
- Card creation, reading, updating, and deletion (CRUD operations) 
- Card attributes for pre-defined types
- Interface-based design for handling cards of different types
- Input validation and error handling
- Query filtering, sorting, and pagination

### Deck Management
- Deck creation and management
- Deck state management (TODO)

### Deck State Management
- Create a mutable instance of an immutable deck (in-progress)
- Shuffle (TODO)
- Draw (TODO)

## Resource Design 

Resources are separated into two main packages, `deckdef` and `deckstate`. 

### Deck Definitions

Package `deckdef` contains definitions for card decks and card types. Game-specific card types may also live in this package, as long as they are purely definitional.

Models in the package are immutable templates that define deck composition and card properties, but do not track gameplay state or card positions. These templates can be used to instantiate one or more independent DeckState instances.

#### Cards 
Cards are structs that implement `CardInterface`. This allows a deck to contain >1 card type, and allows game designers to use multiple card types. 
- `GetID()`
- `GetName()`
- `GetFrontImageURL()`
- `GetBackImageURL()`
- `GetCardType()`

Card Types are extendable. The following are available currently:
- **ImageCard**: Simple cards with just imagery and basic info (name, description, images) ✅
- **PlayingCard**: Standard playing cards (suit, value, images) ✅
- **GameCard**: Generic TCG card for demo (partial implementation)

#### Decks
A deck is an un-ordered collection of structs that implement CardInterface. 

#### Endpoints
- `/v1/image-cards` - Management for the simplest base card type ✅
- `/v1/playing-cards` - Management for traditional playing cards ✅
- `/v1/game-cards` - Generic TCG card management (partial implementation)
- `/v1/decks` - Deck management ✅

### Deck State

Package `deckstate` models the runtime state of a deck during gameplay.

Models in the package track the mutable aspects of a deck instance, such as card zones, card orientation, draw/discard behavior, and per-player hands. This state is created from an immutable deck template (see package deckdef), and evolves as players interact with the cards.

The package makes no assumptions about specific game rules, but provides a flexible framework for modeling card positions, zone behavior, and visibility.

*Implementation is in progress.*

#### DeckState

DeckState is a runtime state of a deck during gameplay. It includes the deck template, player count, and zones where cards are located.

#### Zone

Zone is a collection of cards and groups within a specific area of the game. Zones can represent different game states like draw piles, discard piles, hands, etc.

#### Endpoints
- `/v1/deckstates` - Deck state ✅
- `/v1/deckstates/{id}/actions/{actionName}` - Actions on a deck state (in progress)
  - Supported actions: `add-zone`, `remove-zone`, `sort-zone`, `move-cards`

TODO: Implement move cards. Performance notes:
- Use `slice = slice[:len(slice)-1]` for efficient end-of-slice removal. 
- Indexes are ordered left to right and/or bottom to top
- bottom of deck = start of slice
- left of hand = start of slice
- top of deck = end of slice
- right of hand = end of slice

## Security (TODO)

### Authentication 
1. JWT (JSON Web Tokens) for stateless authentication
2. Rate limiting per API key/user

### Authorization
1. Role-based access control
   - Admin roles for card management
   - Player roles for deck management and simulation