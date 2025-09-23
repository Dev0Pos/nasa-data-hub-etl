package models

import (
	"strconv"
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
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Link        string           `json:"link"`
	Categories  []CategoryObject `json:"categories"`
	Sources     []Source         `json:"sources"`
	Geometry    []Geometry       `json:"geometry"`
	Closed      *string          `json:"closed"`
}

// CategoryObject represents a category object in an event
type CategoryObject struct {
	ID    interface{} `json:"id"` // Can be int or string
	Title string      `json:"title"`
}

// GetIDAsInt converts the ID to int, handling both string and int types
func (co *CategoryObject) GetIDAsInt() int {
	switch v := co.ID.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// Try to parse string as int
		if id, err := strconv.Atoi(v); err == nil {
			return id
		}
		return 0
	default:
		return 0
	}
}

// Category represents an event category
type Category struct {
	ID          interface{} `json:"id"` // Can be int or string from API
	Title       string      `json:"title"`
	Link        string      `json:"link"`
	Description string      `json:"description"`
	Layers      string      `json:"layers"`
}

// GetIDAsInt converts the ID to int, handling both string and int types
func (c *Category) GetIDAsInt() int {
	switch v := c.ID.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// Try to parse string as int
		if id, err := strconv.Atoi(v); err == nil {
			return id
		}
		return 0
	default:
		return 0
	}
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
