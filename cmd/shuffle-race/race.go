package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateDeckStateRequest struct {
	DeckID      string `json:"deck_id"`
	PlayerCount int    `json:"player_count"`
}

type CreateDeckStateResponse struct {
	// Success fields
	ID string `json:"id,omitempty"`

	// Error fields
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

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

	fmt.Printf("✅ Created deck state: %s\n", deckStateID)

	// TODO:
	// Shuffle the deck
	// Draw 5 cards
	// If Royal Flush, end
	// If not, return 5 cards to draw pile
	// Repeat until Royal Flush is drawn

	// Print the number of iterations and time taken

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
