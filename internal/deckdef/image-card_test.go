package deckdef

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestImageCard_Validate(t *testing.T) {
	tests := []struct {
		name    string
		card    ImageCard
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid card",
			card: ImageCard{
				ID:            uuid.New(),
				Name:          "Test Card",
				Description:   "A test card",
				FrontImageURL: "https://example.com/front.jpg",
				BackImageURL:  "https://example.com/back.jpg",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			card: ImageCard{
				Name: "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "whitespace only name",
			card: ImageCard{
				Name: "   ",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "name too long",
			card: ImageCard{
				Name: strings.Repeat("a", 256),
			},
			wantErr: true,
			errMsg:  "name must be 255 characters or less",
		},
		{
			name: "description too long",
			card: ImageCard{
				Name:        "Valid Name",
				Description: strings.Repeat("a", 1001),
			},
			wantErr: true,
			errMsg:  "description must be 1000 characters or less",
		},
		{
			name: "front image URL too long",
			card: ImageCard{
				Name:          "Valid Name",
				FrontImageURL: "https://example.com/" + strings.Repeat("a", 500),
			},
			wantErr: true,
			errMsg:  "front_image_url must be 500 characters or less",
		},
		{
			name: "back image URL too long",
			card: ImageCard{
				Name:         "Valid Name",
				BackImageURL: "https://example.com/" + strings.Repeat("a", 500),
			},
			wantErr: true,
			errMsg:  "back_image_url must be 500 characters or less",
		},
		{
			name: "invalid front image URL",
			card: ImageCard{
				Name:          "Valid Name",
				FrontImageURL: "invalid-url",
			},
			wantErr: true,
			errMsg:  "front_image_url must be a valid URL",
		},
		{
			name: "invalid back image URL",
			card: ImageCard{
				Name:         "Valid Name",
				BackImageURL: "ftp://example.com/image.jpg",
			},
			wantErr: true,
			errMsg:  "back_image_url must be a valid URL",
		},
		{
			name: "empty URLs are valid",
			card: ImageCard{
				Name:          "Valid Name",
				FrontImageURL: "",
				BackImageURL:  "",
			},
			wantErr: false,
		},
		{
			name: "http URLs are valid",
			card: ImageCard{
				Name:          "Valid Name",
				FrontImageURL: "http://example.com/front.jpg",
				BackImageURL:  "http://example.com/back.jpg",
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
