package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nasa-data-hub-etl/internal/config"

	"github.com/sirupsen/logrus"
)

func TestEONETClient_NewEONETClient(t *testing.T) {
	cfg := &config.NASAConfig{
		APIURL: "https://eonet.gsfc.nasa.gov/api/v3",
		APIKey: "test-key",
	}
	logger := logrus.New()

	client := NewEONETClient(cfg, logger)

	if client == nil {
		t.Error("NewEONETClient() returned nil")
		return
	}

	if client.config != cfg {
		t.Error("NewEONETClient() did not set config correctly")
		return
	}

	if client.logger != logger {
		t.Error("NewEONETClient() did not set logger correctly")
		return
	}

	if client.httpClient == nil {
		t.Error("NewEONETClient() did not create HTTP client")
		return
	}
}

func TestEONETClient_FetchCategories(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		statusCode     int
		wantErr        bool
	}{
		{
			name: "successful response - direct array",
			serverResponse: `[
				{
					"id": 8,
					"title": "Wildfires",
					"link": "https://example.com/category/8",
					"description": "Wildfire events",
					"layers": "fire"
				}
			]`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "successful response - object with categories field",
			serverResponse: `{
				"categories": [
					{
						"id": 8,
						"title": "Wildfires",
						"link": "https://example.com/category/8",
						"description": "Wildfire events",
						"layers": "fire"
					}
				]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:           "server error",
			serverResponse: `{"error": "Internal Server Error"}`,
			statusCode:     http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "invalid JSON",
			serverResponse: `invalid json`,
			statusCode:     http.StatusOK,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create client with test server URL
			cfg := &config.NASAConfig{
				APIURL: server.URL,
				APIKey: "",
			}
			logger := logrus.New()
			client := NewEONETClient(cfg, logger)

			ctx := context.Background()
			categories, err := client.FetchCategories(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if categories == nil {
					t.Error("FetchCategories() returned nil categories")
					return
				}

				if len(categories) == 0 {
					t.Error("FetchCategories() returned empty categories")
					return
				}

				// Validate first category
				category := categories[0]
				if category.Title == "" {
					t.Error("Category Title is empty")
				}

				// Test GetIDAsInt method
				id := category.GetIDAsInt()
				if id == 0 {
					t.Error("Category ID is zero")
				}
			}
		})
	}
}

func TestEONETClient_ContextCancellation(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	// Create client
	cfg := &config.NASAConfig{
		APIURL: server.URL,
		APIKey: "",
	}
	logger := logrus.New()
	client := NewEONETClient(cfg, logger)

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.FetchCategories(ctx)

	if err == nil {
		t.Error("FetchCategories() should have returned error due to context cancellation")
	}
}
