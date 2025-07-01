package deckdef

import (
	"testing"

	"github.com/google/uuid"
)

func TestDeckInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		deck    DeckInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid deck",
			deck: DeckInput{
				Name: "Test Deck",
				Type: "standard",
				Cards: &CardCollectionInput{
					Items: []CardInputWithQuantity{
						{
							Card:     CardInput{ID: uuid.New()},
							Quantity: 2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid deck without cards",
			deck: DeckInput{
				Name: "Test Deck",
				Type: "standard",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			deck: DeckInput{
				Name: "",
				Type: "standard",
			},
			wantErr: true,
			errMsg:  "deck name is required",
		},
		{
			name: "nil card ID",
			deck: DeckInput{
				Name: "Test Deck",
				Type: "standard",
				Cards: &CardCollectionInput{
					Items: []CardInputWithQuantity{
						{
							Card:     CardInput{ID: uuid.Nil},
							Quantity: 1,
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "card ID is required for cards.items[0]",
		},
		{
			name: "zero quantity",
			deck: DeckInput{
				Name: "Test Deck",
				Type: "standard",
				Cards: &CardCollectionInput{
					Items: []CardInputWithQuantity{
						{
							Card:     CardInput{ID: uuid.New()},
							Quantity: 0,
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "card quantity must be positive for cards.items[0]",
		},
		{
			name: "negative quantity",
			deck: DeckInput{
				Name: "Test Deck",
				Type: "standard",
				Cards: &CardCollectionInput{
					Items: []CardInputWithQuantity{
						{
							Card:     CardInput{ID: uuid.New()},
							Quantity: -1,
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "card quantity must be positive for cards.items[0]",
		},
		{
			name: "multiple cards with error on second",
			deck: DeckInput{
				Name: "Test Deck",
				Type: "standard",
				Cards: &CardCollectionInput{
					Items: []CardInputWithQuantity{
						{
							Card:     CardInput{ID: uuid.New()},
							Quantity: 1,
						},
						{
							Card:     CardInput{ID: uuid.Nil},
							Quantity: 1,
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "card ID is required for cards.items[1]",
		},
		{
			name: "valid deck with sleeve image URL",
			deck: DeckInput{
				Name:           "Test Deck",
				Type:           "standard",
				SleeveImageURL: stringPtr("https://example.com/sleeve.jpg"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.deck.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
