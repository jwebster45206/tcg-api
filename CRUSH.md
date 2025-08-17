# CRUSH.md

Project: tcg-api (Go)

## Commands
- Build: go build ./...
- Test (unit): go test ./...
- Test (verbose,race,cover): go test -v -race -coverprofile=coverage.out ./...
- Lint (Go): golangci-lint run
- Build Docker: docker build -t tcg-api:local .
- Run Docker Compose: docker compose up --build

## Auth
- Issue a token (CLI): go run ./cmd/jwt-issuer --sub admin --scopes admin --ttl 1h
- Use token: add header Authorization: Bearer <token>
- Health endpoint is open; all others are protected when auth.enabled=true

## Entrypoints
- cmd/tcg-api/main.go — REST API server
- cmd/shuffle-race/main.go — shuffle benchmarking tool
- cmd/jwt-issuer/main.go — dev CLI to mint JWTs

## Key Packages
- internal/auth — JWT claims, issuer/verifier (HS256), middleware
- internal/handlers — HTTP handlers for resources (health, cards, decks, deckstate, shuffle, zones)
- internal/storage — persistence layer (MySQL)
- internal/state — deckstate storage abstraction
- internal/config — config loading and logging
- internal/shuffle — shuffle strategies and tests
- internal/query — query parsing
- pkg/deckdef — immutable deck and card definitions (image, playing, game)
- pkg/deckstate — runtime deck state, zones, operations

## Configuration
- config.json.sample — example app config (includes auth block)
- config.docker.json — defaults for Docker (includes auth block)
- Environment (Dockerfile): PORT=8080
- DB: MySQL (see db/schema.sql and seeds.sql). Redis noted in README but not required at runtime for tests.

## Make it easy next time
I’ll use the commands above by default. If you add Makefile or scripts later, I’ll mirror them here.
