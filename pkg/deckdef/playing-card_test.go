package deckdef

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPlayingCard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    PlayingCard
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid card - ace of hearts",
			card: PlayingCard{
				ID:            uuid.New(),
				Name:          "Ace of Hearts",
				Suit:          SuitHearts,
				Ranking:       1,
				FrontImageURL: "https://example.com/front.jpg",
				BackImageURL:  "https://example.com/back.jpg",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid card - king of spades",
			card: PlayingCard{
				Suit:    SuitSpades,
				Ranking: 13,
			},
			wantErr: false,
		},
		{
			name: "ranking too low",
			card: PlayingCard{
				Suit:    SuitHearts,
				Ranking: 0,
			},
			wantErr: true,
			errMsg:  "ranking must be between 1 and 13",
		},
		{
			name: "ranking too high",
			card: PlayingCard{
				Suit:    SuitDiamonds,
				Ranking: 14,
			},
			wantErr: true,
			errMsg:  "ranking must be between 1 and 13",
		},
		{
			name: "invalid suit",
			card: PlayingCard{
				Suit:    "jokers",
				Ranking: 1,
			},
			wantErr: true,
			errMsg:  "invalid suit: jokers",
		},
		{
			name: "empty suit",
			card: PlayingCard{
				Suit:    "",
				Ranking: 1,
			},
			wantErr: true,
			errMsg:  "invalid suit: ",
		},
		{
			name: "valid suits test - hearts",
			card: PlayingCard{
				Suit:    SuitHearts,
				Ranking: 7,
			},
			wantErr: false,
		},
		{
			name: "valid suits test - diamonds",
			card: PlayingCard{
				Suit:    SuitDiamonds,
				Ranking: 7,
			},
			wantErr: false,
		},
		{
			name: "valid suits test - clubs",
			card: PlayingCard{
				Suit:    SuitClubs,
				Ranking: 7,
			},
			wantErr: false,
		},
		{
			name: "valid suits test - spades",
			card: PlayingCard{
				Suit:    SuitSpades,
				Ranking: 7,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate()
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
