package gojpcal

import (
	"time"
)

// Event represents a connpass event
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"started_at"`
	EndTime     time.Time `json:"ended_at"`
	Place       string    `json:"place"`
	Address     string    `json:"address"`
}

// ConnpassResponse represents the response from the Connpass API
type ConnpassResponse struct {
	ResultsReturned  int     `json:"results_returned"`
	ResultsStart     int     `json:"results_start"`
	ResultsAvailable int     `json:"results_available"`
	Events           []Event `json:"events"`
}
