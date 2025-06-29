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

## Resource Design

### Cards 
Cards are structs that implement CardInterface. This allows a deck to contain >1 card type, and allows game designers to use multiple card types. 

#### CardInterface**
Base contract that all card types implement
- `GetID()`
- `GetName()`
- `GetFrontImageURL()`
- `GetBackImageURL()`
- `GetCardType()`

#### Implementations
- **ImageCard**: Simple cards with just imagery and basic info (name, description, images) ✅
- **PlayingCard**: Standard playing cards (suit, value, images) ✅
- **GameCard**: TCG-specific cards with game mechanics (cost, offense, defense, keywords, colors) - *Partial*

### Decks
A deck is an ordered collection of structs that implement a basic card interface. Deck and card interfaces should be agnostic of game mechanics. If a mechanic is specific to a game and not widely applicable to most decks of cards, it does not belong in the deck api. 

## API Endpoints
- `/v1/image-cards` - Management for the simplest base card type ✅
- `/v1/playing-cards` - Management for traditional playing cards ✅
- `/v1/game-cards` - TCG card management - partial/demo implementation
- `/v1/decks` - Deck management ✅
- Deck state / deck operations: TODO

## Security (TODO)

### Authentication 
1. JWT (JSON Web Tokens) for stateless authentication
2. Rate limiting per API key/user

### Authorization
1. Role-based access control
   - Admin roles for card management
   - Player roles for deck management and simulation