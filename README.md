# tcg-api
Golang API for trading card game deck building

## Overview
A lightweight REST API for simulating trading card games, built with Go and following HAL (Hypertext Application Language) standards.

## Technical Stack
- **Language**: Go
- **Containerization**: Docker
- **API Style**: REST with HAL
- **Storage**: MySQL and Redis

## Core Features

### Card Management
- Card creation and management
- Card attributes (name, cost, offense, defense, keywords, colors)

### Card Deck Management
- Deck creation and management
- Validation of decks
- Statistics of decks

### Game Mechanics
- Simulation of shuffle and draw
- Other mechanics TBD

### API Endpoints
- `/game-cards` - Card resource management
- `/decks` - Deck management
- TODO - shuffle and draw

## Storage Architecture

### Primary Storage: MySQL
- Stores persistent data (cards, decks, users, game history)
- Handles complex relationships and transactions
- ACID compliance for data integrity
- Good performance with proper indexing

### Cache Layer: Redis
- Fast access to frequently queried data
- Session storage
- Game state caching
- Rate limiting counters
- Deck validation caching

## Infra

### Local Development
- Docker Compose stack with MySQL and Redis containers
- Volume mounting for data persistence during development
- Environment-based configuration

### Production Ready
- Designed for cloud deployment (EKS + RDS + ElastiCache)
- Connection pooling and retry logic
- Health checks and monitoring endpoints
- Graceful shutdown handling

## Security

### Authentication
1. JWT (JSON Web Tokens) for stateless authentication
2. Rate limiting per API key/user

### Authorization
1. Role-based access control (RBAC)
   - Admin roles for card management
   - Player roles for deck management and simulation

### API Security
1. TLS encryption (mandatory)
2. API key management
3. Input validation and sanitization
4. CORS policy implementation

## Development Phases

### Phase 1: Foundation
1. Project setup with Go modules
2. Basic API structure with HAL implementation
3. Docker configuration
4. CI/CD pipeline setup
5. Storage layer implementation

### Phase 2: Deck Management
1. Card management system
2. Deck building functionality

### Phase 3: Deck Simulation
1. TBD