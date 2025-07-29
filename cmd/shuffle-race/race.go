package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jwebster45206/tcg-api/pkg/deckdef"
	"github.com/jwebster45206/tcg-api/pkg/deckstate"
)

const maxTries = 1000

func race(decks int) error {
	if decks != 1 {
		return fmt.Errorf("only 1 deck supported for now")
	}

	client := newHTTPClient()

	fmt.Println("Starting shuffle race with", decks, "deck(s)...")

	// Create deckstate for standard playing card deck
	deckStateID, err := createDeckState(client)
	if err != nil {
		return fmt.Errorf("failed to create deck state: %w", err)
	}

	fmt.Printf("Created deck state: %s\n", deckStateID)

	// Try up to maxTries times to find a Royal Flush
	for iteration := 1; iteration <= maxTries; iteration++ {
		fmt.Printf("\n--- Iteration %d ---\n", iteration)

		// Shuffle the deck
		err = shuffleDeck(client, deckStateID)
		if err != nil {
			return fmt.Errorf("failed to shuffle deck: %w", err)
		}
		fmt.Println("Shuffled the deck")

		// Draw 5 cards
		err = drawCards(client, deckStateID, "draw", "player:1", 5)
		if err != nil {
			return fmt.Errorf("failed to draw cards: %w", err)
		}

		fmt.Println("Drew 5 cards to player:1")

		deckState, err := getDeckState(client, deckStateID)
		if err != nil {
			return fmt.Errorf("failed to get deck state: %w", err)
		}
		printPlayerHand(deckState)

		// Check if Royal Flush
		isRoyalFlush := checkRoyalFlush(deckState)
		if isRoyalFlush {
			fmt.Printf("ROYAL FLUSH FOUND after %d iterations!\n", iteration)
			return nil
		}

		fmt.Println("No Royal Flush found")

		// Return 5 cards to draw pile for next iteration
		err = drawCards(client, deckStateID, "player:1", "draw", 5)
		if err != nil {
			return fmt.Errorf("failed to return cards to draw pile: %w", err)
		}
		fmt.Println("Returned 5 cards to draw pile")
	}

	fmt.Println("No Royal Flush found after 100 iterations")

	return nil
}

// getDeckState retrieves a deck state by ID using the API
func getDeckState(client *http.Client, deckStateID string) (*deckstate.DeckState, error) {
	// Create HTTP request
	apiURL := APIBaseURL + "/v1/deckstates/" + deckStateID
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Success case - parse as DeckState
		var deckState deckstate.DeckState
		if err := json.NewDecoder(resp.Body).Decode(&deckState); err != nil {
			return nil, fmt.Errorf("failed to parse deck state response: %w", err)
		}
		return &deckState, nil
	} else {
		// Error case - parse as ErrorResponse
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("API returned status %d (could not parse error response: %w)", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, errorResponse.Error, errorResponse.Message)
	}
}

// drawCards draws cards from one zone to another using the API
func drawCards(client *http.Client, deckStateID, fromZone, toZone string, count int) error {
	request := DrawRequest{
		FromZone: fromZone,
		ToZone:   toZone,
		Count:    count,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := APIBaseURL + "/v1/deckstates/" + deckStateID + "/actions/draw-cards"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var response DrawResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if response.Error != "" {
			return fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, response.Error, response.Message)
		}
		return fmt.Errorf("API returned status %d with no error details", resp.StatusCode)
	}

	if !response.Success {
		return fmt.Errorf("draw operation was not successful")
	}

	return nil
}

// shuffleDeck shuffles the draw zone of a deck state using the API
func shuffleDeck(client *http.Client, deckStateID string) error {
	request := SortZoneRequest{
		Zone: "draw",
		Sort: "shuffle",
	}

	fmt.Println("Shuffling deck...")
	fmt.Printf("  Deck State ID: %s\n", deckStateID)
	fmt.Printf("  Zone: %s\n", request.Zone)
	fmt.Printf("  Sort Type: %s\n", request.Sort)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := APIBaseURL + "/v1/deckstates/" + deckStateID + "/actions/sort-zone"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var response SortZoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if response.Error != "" {
			return fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, response.Error, response.Message)
		}
		return fmt.Errorf("API returned status %d with no error details", resp.StatusCode)
	}

	if !response.Success {
		return fmt.Errorf("shuffle operation was not successful")
	}

	return nil
}

// createDeckState creates a new deck state using the API
func createDeckState(client *http.Client) (string, error) {
	request := CreateDeckStateRequest{
		DeckID:      "d0000000-0000-0000-0000-000000000001", // Standard playing card deck
		PlayerCount: 1,
	}

	fmt.Println("Creating deck state...")
	fmt.Printf("  Deck ID: %s\n", request.DeckID)
	fmt.Printf("  Player Count: %d\n", request.PlayerCount)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := APIBaseURL + "/v1/deckstates"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var response CreateDeckStateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		if response.Error != "" {
			return "", fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, response.Error, response.Message)
		}
		return "", fmt.Errorf("API returned status %d with no error details", resp.StatusCode)
	}

	if response.ID == "" {
		return "", fmt.Errorf("API response missing deck state ID")
	}

	return response.ID, nil
}

// printPlayerHand prints the cards in the player's hand
func printPlayerHand(deckState *deckstate.DeckState) {
	playingCards := extractPlayingCards(deckState, "player:1")
	fmt.Printf("Playing Cards in player's hand (%d):\n", len(playingCards))
	for i, card := range playingCards {
		fmt.Printf("  %d. %s\n", i+1, card.GetName())
	}
}

var royalFlushRankings = map[int]bool{1: true, 10: true, 11: true, 12: true, 13: true}

// checkRoyalFlush checks if the "player:1" zone contains a Royal Flush
// Where a Royal Flush is: 10, Jack, Queen, King, Ace of the same suit
func checkRoyalFlush(deckState *deckstate.DeckState) bool {
	playingCards := extractPlayingCards(deckState, "player:1")
	if len(playingCards) < 5 {
		return false
	}

	// Group cards by suit
	suitCards := make(map[string][]int) // suite -> rankings
	for _, card := range playingCards {
		suitCards[card.Suit] = append(suitCards[card.Suit], card.Ranking)
	}
	if len(suitCards) > 1 {
		return false
	}

	cards := suitCards[playingCards[0].Suit]
	rankingsMap := make(map[int]bool)
	for _, ranking := range cards {
		rankingsMap[ranking] = true
	}

	// Check if all Royal Flush rankings are present
	hasRoyalFlush := true
	for requiredRanking := range royalFlushRankings {
		if !rankingsMap[requiredRanking] {
			hasRoyalFlush = false
			break
		}
	}
	return hasRoyalFlush
}

// extractPlayingCards extracts playingCard data from a zone
func extractPlayingCards(deckState *deckstate.DeckState, zoneName string) []*deckdef.PlayingCard {
	zone, exists := deckState.Zones[zoneName]
	if !exists || zone == nil {
		return nil
	}

	var playingCards []*deckdef.PlayingCard
	for _, item := range zone.Items {
		if cardInZone, ok := item.(deckstate.CardInZone); ok {
			if playingCard, ok := cardInZone.Card.(*deckdef.PlayingCard); ok {
				playingCards = append(playingCards, playingCard)
			}
		}
	}
	return playingCards
}
