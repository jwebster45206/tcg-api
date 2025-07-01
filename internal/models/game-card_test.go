package models

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGameCard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    GameCard
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid card",
			card: GameCard{
				ID:            uuid.New(),
				Name:          "Lightning Bolt",
				Description:   "Deal 3 damage to any target",
				ManaCost:      1,
				CardType:      "instant",
				Attack:        0,
				Health:        0,
				Keywords:      []string{"instant"},
				Colors:        []string{"red"},
				Rarity:        "common",
				SetCode:       "TST",
				FrontImageURL: "https://example.com/front.jpg",
				BackImageURL:  "https://example.com/back.jpg",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			card: GameCard{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "whitespace only name",
			card: GameCard{
				Name: "   ",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too long",
			card: GameCard{
				Name: strings.Repeat("a", 256),
			},
			wantErr: true,
			errMsg:  "name must be 255 characters or less",
		},
		{
			name: "description too long",
			card: GameCard{
				Name:        "Valid Name",
				Description: strings.Repeat("a", 1001),
			},
			wantErr: true,
			errMsg:  "description must be 1000 characters or less",
		},
		{
			name: "negative mana cost",
			card: GameCard{
				Name:     "Valid Name",
				ManaCost: -1,
			},
			wantErr: true,
			errMsg:  "mana_cost cannot be negative",
		},
		{
			name: "negative attack",
			card: GameCard{
				Name:   "Valid Name",
				Attack: -1,
			},
			wantErr: true,
			errMsg:  "attack cannot be negative",
		},
		{
			name: "negative health",
			card: GameCard{
				Name:   "Valid Name",
				Health: -1,
			},
			wantErr: true,
			errMsg:  "health cannot be negative",
		},
		{
			name: "invalid rarity",
			card: GameCard{
				Name:   "Valid Name",
				Rarity: "super-rare",
			},
			wantErr: true,
			errMsg:  "rarity must be one of: common, uncommon, rare, mythic",
		},
		{
			name: "set code too long",
			card: GameCard{
				Name:    "Valid Name",
				SetCode: "TOOLONGCODE",
			},
			wantErr: true,
			errMsg:  "set_code must be 10 characters or less",
		},
		{
			name: "valid rarities",
			card: GameCard{
				Name:   "Valid Name",
				Rarity: "legendary",
			},
			wantErr: false,
		},
		{
			name: "empty rarity is valid",
			card: GameCard{
				Name:   "Valid Name",
				Rarity: "",
			},
			wantErr: false,
		},
		{
			name: "zero values are valid",
			card: GameCard{
				Name:     "Valid Name",
				ManaCost: 0,
				Attack:   0,
				Health:   0,
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
