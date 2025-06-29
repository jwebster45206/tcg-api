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

// CardWithQuantityJSON allows marshaling polymorphic card types with quantity.
type CardWithQuantityJSON struct {
	Card     json.RawMessage `json:"card"`
	Quantity int             `json:"quantity"`
}

func (cwq CardWithQuantity) MarshalJSON() ([]byte, error) {
	// Marshal the card directly, then add card_type
	cardJSON, err := json.Marshal(cwq.Card)
	if err != nil {
		return nil, err
	}

	// Parse the marshaled card to add card_type
	var cardMap map[string]interface{}
	if err := json.Unmarshal(cardJSON, &cardMap); err != nil {
		return nil, err
	}

	// Add the card_type field
	cardMap["card_type"] = cwq.Card.GetCardType()

	// Re-marshal with card_type included
	cardWithTypeJSON, err := json.Marshal(cardMap)
	if err != nil {
		return nil, err
	}

	result := struct {
		Card     json.RawMessage `json:"card"`
		Quantity int             `json:"quantity"`
	}{
		Card:     cardWithTypeJSON,
		Quantity: cwq.Quantity,
	}

	return json.Marshal(result)
}

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
