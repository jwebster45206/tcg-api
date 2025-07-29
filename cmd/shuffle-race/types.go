package main

type CreateDeckStateRequest struct {
	DeckID      string `json:"deck_id"`
	PlayerCount int    `json:"player_count"`
}

type CreateDeckStateResponse struct {
	ID      string `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type SortZoneRequest struct {
	Zone string `json:"zone"`
	Sort string `json:"sort"`
}

type SortZoneResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type DrawRequest struct {
	FromZone string `json:"from_zone"`
	ToZone   string `json:"to_zone"`
	Count    int    `json:"count"`
}

type DrawResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
