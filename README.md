# tcg-api
A lightweight REST API for simulating card decks, built with Go. The API features a modular card interface design that separates game-specific mechanics from general deck mechanics, making it extensible for different types of card games.

## Technical Stack
- **Language**: Go
- **Storage**: MySQL and Redis (TODO)

## Core Features

### Card Management
- Card creation and management (in-progress)
- Card attributes for pre-defined types: image-card, playing-card, and game-card (in-progress)
- Interface-based design for handling cards of different types

### Deck Management
- Deck creation and management (TODO)
- Deck state management (TODO)

## Architecture Design

### Card Interface System
The API uses an interface-driven approach to support multiple card types:

- **CardInterface**: Base contract that all card types implement
  - `GetID()`, `GetName()`, `GetFrontImageURL()`, `GetBackImageURL()`, `GetCardType()`

### Card Type Implementations
- **GameCard**: TCG-specific cards with game mechanics (cost, offense, defense, keywords, colors)
- **ImageCard**: Simple cards with just imagery and basic info (name, description, images)
- **PlayingCard**: Standard playing cards (suite, value, images) - *TODO*

### Deck Type Implementation (TODO)
- Array of cards (unsorted) by identifier and type
- Name of deck
- Owner of deck (can be nil)
- Sleeve/Back image URL (can be nil)
- Deck accepts cards of any type implementing CardInterface

## API Endpoints
- `/game-cards` - GameCard resource management (TCG-specific cards)
- `/decks` - Deck management (TODO)
- TODO - ImageCard and PlayingCard handlers
- TODO - shuffle and draw

## Security

### Authentication (TODO)
1. JWT (JSON Web Tokens) for stateless authentication
2. Rate limiting per API key/user

### Authorization
1. Role-based access control (RBAC)
   - Admin roles for card management
   - Player roles for deck management and simulation