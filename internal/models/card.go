package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// CardInterface defines the contract that all card types must implement
type CardInterface interface {
	GetID() uuid.UUID
	GetName() string
	GetFrontImageURL() string
	GetBackImageURL() string
	GetCardType() string // Used for routing to correct storage/handlers
}

// CardWithQuantityJSON represents the JSON structure for CardWithQuantity
type CardWithQuantityJSON struct {
	Card     json.RawMessage `json:"card"`
	Quantity int             `json:"quantity"`
}

// MarshalJSON implements custom JSON marshaling for CardWithQuantity
func (cwq CardWithQuantity) MarshalJSON() ([]byte, error) {
	// Create a wrapper struct with the card type information
	cardData := struct {
		CardInterface
		CardType string `json:"card_type"`
	}{
		CardInterface: cwq.Card,
		CardType:      cwq.Card.GetCardType(),
	}

	cardJSON, err := json.Marshal(cardData)
	if err != nil {
		return nil, err
	}

	result := struct {
		Card     json.RawMessage `json:"card"`
		Quantity int             `json:"quantity"`
	}{
		Card:     cardJSON,
		Quantity: cwq.Quantity,
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for CardWithQuantity
func (cwq *CardWithQuantity) UnmarshalJSON(data []byte) error {
	var temp CardWithQuantityJSON
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	cwq.Quantity = temp.Quantity

	// Parse the card type to determine which struct to unmarshal into
	var cardTypeInfo struct {
		CardType string `json:"card_type"`
	}
	if err := json.Unmarshal(temp.Card, &cardTypeInfo); err != nil {
		return err
	}

	// Unmarshal into the appropriate card type
	switch cardTypeInfo.CardType {
	case TypePlayingCard:
		var card PlayingCard
		if err := json.Unmarshal(temp.Card, &card); err != nil {
			return err
		}
		cwq.Card = &card
	case TypeImageCard:
		var card ImageCard
		if err := json.Unmarshal(temp.Card, &card); err != nil {
			return err
		}
		cwq.Card = &card
	case TypeGameCard:
		var card GameCard
		if err := json.Unmarshal(temp.Card, &card); err != nil {
			return err
		}
		cwq.Card = &card
	default:
		return fmt.Errorf("unknown card type: %s", cardTypeInfo.CardType)
	}

	return nil
}
