package models

import (
	"time"
)

// EONETResponse represents the response from NASA EONET API
type EONETResponse struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	Events      []Event    `json:"events"`
	Categories  []Category `json:"categories"`
}

// Event represents a natural event from EONET
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Categories  []int     `json:"categories"`
	Sources     []Source  `json:"sources"`
	Geometry    []Geometry `json:"geometry"`
	Closed      *string   `json:"closed"`
}

// Category represents an event category
type Category struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Layers      string `json:"layers"`
}

// Source represents a data source for an event
type Source struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// Geometry represents the geographic data for an event
type Geometry struct {
	Date        time.Time   `json:"date"`
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

// EventRecord represents a processed event record for database storage
type EventRecord struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Link        string    `db:"link"`
	Categories  string    `db:"categories"` // JSON string
	Sources     string    `db:"sources"`    // JSON string
	Geometry    string    `db:"geometry"`   // JSON string
	Closed      *string   `db:"closed"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// CategoryRecord represents a category record for database storage
type CategoryRecord struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	Description string    `db:"description"`
	Layers      string    `db:"layers"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
