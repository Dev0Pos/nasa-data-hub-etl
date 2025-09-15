package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"nasa-data-hub-etl/internal/config"
	"nasa-data-hub-etl/pkg/models"

	"github.com/sirupsen/logrus"
)

// EONETClient handles communication with NASA EONET API
type EONETClient struct {
	config     *config.NASAConfig
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewEONETClient creates a new EONET API client
func NewEONETClient(cfg *config.NASAConfig, logger *logrus.Logger) *EONETClient {
	return &EONETClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// FetchEventsOptions represents options for fetching events
type FetchEventsOptions struct {
	Days       int    `json:"days,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Status     string `json:"status,omitempty"` // "open", "closed", "all"
	CategoryID int    `json:"category,omitempty"`
	SourceID   string `json:"source,omitempty"`
}

// FetchEvents fetches events from NASA EONET API
func (c *EONETClient) FetchEvents(ctx context.Context, opts FetchEventsOptions) (*models.EONETResponse, error) {
	url := c.buildEventsURL(opts)

	c.logger.WithFields(logrus.Fields{
		"url":  url,
		"opts": opts,
	}).Debug("Fetching events from NASA EONET API")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if provided
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	req.Header.Set("User-Agent", "NASA-Data-Hub-ETL/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var eonetResponse models.EONETResponse
	if err := json.Unmarshal(body, &eonetResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"events_count":     len(eonetResponse.Events),
		"categories_count": len(eonetResponse.Categories),
	}).Info("Successfully fetched events from NASA EONET API")

	return &eonetResponse, nil
}

// FetchCategories fetches categories from NASA EONET API
func (c *EONETClient) FetchCategories(ctx context.Context) ([]models.Category, error) {
	url := fmt.Sprintf("%s/categories", c.config.APIURL)

	c.logger.WithField("url", url).Debug("Fetching categories from NASA EONET API")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if provided
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	req.Header.Set("User-Agent", "NASA-Data-Hub-ETL/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var categories []models.Category
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.WithField("categories_count", len(categories)).Info("Successfully fetched categories from NASA EONET API")

	return categories, nil
}

// buildEventsURL builds the URL for fetching events with the given options
func (c *EONETClient) buildEventsURL(opts FetchEventsOptions) string {
	url := fmt.Sprintf("%s/events", c.config.APIURL)

	params := make([]string, 0)

	if opts.Days > 0 {
		params = append(params, fmt.Sprintf("days=%d", opts.Days))
	}

	if opts.Limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", opts.Limit))
	}

	if opts.Status != "" {
		params = append(params, fmt.Sprintf("status=%s", opts.Status))
	}

	if opts.CategoryID > 0 {
		params = append(params, fmt.Sprintf("category=%d", opts.CategoryID))
	}

	if opts.SourceID != "" {
		params = append(params, fmt.Sprintf("source=%s", opts.SourceID))
	}

	if len(params) > 0 {
		url += "?" + params[0]
		for i := 1; i < len(params); i++ {
			url += "&" + params[i]
		}
	}

	return url
}

// HealthCheck checks if the NASA EONET API is accessible
func (c *EONETClient) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/categories", c.config.APIURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("User-Agent", "NASA-Data-Hub-ETL/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}
