package models

import (
	"testing"
	"time"
)

func TestCategory_GetIDAsInt(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		expected int
	}{
		{
			name: "int ID",
			category: Category{
				ID:    8,
				Title: "Wildfires",
			},
			expected: 8,
		},
		{
			name: "float64 ID",
			category: Category{
				ID:    8.0,
				Title: "Wildfires",
			},
			expected: 8,
		},
		{
			name: "string ID",
			category: Category{
				ID:    "8",
				Title: "Wildfires",
			},
			expected: 8,
		},
		{
			name: "invalid string ID",
			category: Category{
				ID:    "invalid",
				Title: "Wildfires",
			},
			expected: 0,
		},
		{
			name: "nil ID",
			category: Category{
				ID:    nil,
				Title: "Wildfires",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.category.GetIDAsInt()
			if result != tt.expected {
				t.Errorf("GetIDAsInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCategoryObject_GetIDAsInt(t *testing.T) {
	tests := []struct {
		name     string
		category CategoryObject
		expected int
	}{
		{
			name: "int ID",
			category: CategoryObject{
				ID:    8,
				Title: "Wildfires",
			},
			expected: 8,
		},
		{
			name: "float64 ID",
			category: CategoryObject{
				ID:    8.0,
				Title: "Wildfires",
			},
			expected: 8,
		},
		{
			name: "string ID",
			category: CategoryObject{
				ID:    "8",
				Title: "Wildfires",
			},
			expected: 8,
		},
		{
			name: "invalid string ID",
			category: CategoryObject{
				ID:    "invalid",
				Title: "Wildfires",
			},
			expected: 0,
		},
		{
			name: "nil ID",
			category: CategoryObject{
				ID:    nil,
				Title: "Wildfires",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.category.GetIDAsInt()
			if result != tt.expected {
				t.Errorf("GetIDAsInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEventRecord_Validation(t *testing.T) {
	tests := []struct {
		name    string
		record  EventRecord
		wantErr bool
	}{
		{
			name: "valid record",
			record: EventRecord{
				ID:          "EONET_12345",
				Title:       "Test Event",
				Description: "Test Description",
				Link:        "https://example.com",
				Categories:  `[8, 12]`,
				Sources:     `[{"id":"test","url":"https://test.com","title":"Test Source"}]`,
				Geometry:    `[{"date":"2025-01-22T10:30:00Z","type":"Point","coordinates":[-120.5,37.8]}]`,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			record: EventRecord{
				ID:          "",
				Title:       "Test Event",
				Description: "Test Description",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty title",
			record: EventRecord{
				ID:          "EONET_12345",
				Title:       "",
				Description: "Test Description",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEventRecord(tt.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEventRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryRecord_Validation(t *testing.T) {
	tests := []struct {
		name    string
		record  CategoryRecord
		wantErr bool
	}{
		{
			name: "valid record",
			record: CategoryRecord{
				ID:          8,
				Title:       "Wildfires",
				Link:        "https://example.com",
				Description: "Wildfire events",
				Layers:      "fire",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "zero ID",
			record: CategoryRecord{
				ID:        0,
				Title:     "Wildfires",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty title",
			record: CategoryRecord{
				ID:        8,
				Title:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCategoryRecord(tt.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCategoryRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper validation functions
func validateEventRecord(record EventRecord) error {
	if record.ID == "" {
		return &ValidationError{Field: "ID", Message: "ID cannot be empty"}
	}
	if record.Title == "" {
		return &ValidationError{Field: "Title", Message: "Title cannot be empty"}
	}
	return nil
}

func validateCategoryRecord(record CategoryRecord) error {
	if record.ID == 0 {
		return &ValidationError{Field: "ID", Message: "ID cannot be zero"}
	}
	if record.Title == "" {
		return &ValidationError{Field: "Title", Message: "Title cannot be empty"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
