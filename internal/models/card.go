package models

import "github.com/google/uuid"

type GameCard struct {
	ID         uuid.UUID
	Name       string   `json:"name"`
	Subtitle   string   `json:"subtitle"`
	Cost       int      `json:"cost"`
	Type       string   `json:"type"`
	Offense    int      `json:"offense"`
	Defense    int      `json:"defense"`
	Keywords   []string `json:"keywords"`
	Colors     []string `json:"colors"`
	IsResource bool     `json:"is_resource"`
}
