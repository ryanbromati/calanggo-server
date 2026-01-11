package domain

import "time"

type Link struct {
	ID        string    `json:"id"`
	Original  string    `json:"original"`
	Shortened string    `json:"shortened"`
	CreatedAt time.Time `json:"created_at"`
	Visits    int       `json:"visits"`
}

func NewLink(original string, shortCode string) *Link {
	return &Link{
		ID:        shortCode,
		Original:  original,
		Shortened: shortCode,
		CreatedAt: time.Now(),
		Visits:    0,
	}
}
