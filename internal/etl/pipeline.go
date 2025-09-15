package etl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nasa-data-hub-etl/internal/api"
	"nasa-data-hub-etl/internal/config"
	"nasa-data-hub-etl/internal/database"
	"nasa-data-hub-etl/pkg/models"

	"github.com/sirupsen/logrus"
)

// Pipeline represents the ETL pipeline
type Pipeline struct {
	config     *config.Config
	eonetClient *api.EONETClient
	db         *database.VerticaDB
	logger     *logrus.Logger
}

// NewPipeline creates a new ETL pipeline
func NewPipeline(cfg *config.Config, logger *logrus.Logger) (*Pipeline, error) {
	// Create NASA EONET API client
	eonetClient := api.NewEONETClient(&cfg.NASA, logger)

	// Create VerticaDB connection
	db, err := database.NewVerticaDB(&cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return &Pipeline{
		config:     cfg,
		eonetClient: eonetClient,
		db:         db,
		logger:     logger,
	}, nil
}

// Run starts the ETL pipeline
func (p *Pipeline) Run(ctx context.Context) error {
	p.logger.Info("Starting ETL pipeline")

	// Start ETL run tracking
	runID, err := p.db.StartETLRun(ctx)
	if err != nil {
		return fmt.Errorf("failed to start ETL run tracking: %w", err)
	}

	var eventsProcessed, categoriesProcessed int
	var finalError error

	defer func() {
		status := "completed"
		var errorMsg *string
		if finalError != nil {
			status = "failed"
			msg := finalError.Error()
			errorMsg = &msg
		}

		if err := p.db.CompleteETLRun(ctx, runID, status, eventsProcessed, categoriesProcessed, errorMsg); err != nil {
			p.logger.WithError(err).Error("Failed to complete ETL run tracking")
		}
	}()

	// Process categories first
	if err := p.processCategories(ctx); err != nil {
		finalError = fmt.Errorf("failed to process categories: %w", err)
		return finalError
	}
	categoriesProcessed = 1 // We process all categories in one batch

	// Process events
	eventsProcessed, err = p.processEvents(ctx)
	if err != nil {
		finalError = fmt.Errorf("failed to process events: %w", err)
		return finalError
	}

	p.logger.WithFields(logrus.Fields{
		"events_processed":     eventsProcessed,
		"categories_processed": categoriesProcessed,
	}).Info("ETL pipeline completed successfully")

	return nil
}

// processCategories fetches and processes categories
func (p *Pipeline) processCategories(ctx context.Context) error {
	p.logger.Info("Processing categories")

	// Fetch categories from NASA EONET API
	categories, err := p.eonetClient.FetchCategories(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %w", err)
	}

	// Transform categories to database records
	categoryRecords := make([]*models.CategoryRecord, 0, len(categories))
	for _, category := range categories {
		record := &models.CategoryRecord{
			ID:          category.ID,
			Title:       category.Title,
			Link:        category.Link,
			Description: category.Description,
			Layers:      category.Layers,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		categoryRecords = append(categoryRecords, record)
	}

	// Batch insert categories
	if err := p.db.BatchInsertCategories(ctx, categoryRecords); err != nil {
		return fmt.Errorf("failed to insert categories: %w", err)
	}

	p.logger.WithField("count", len(categoryRecords)).Info("Successfully processed categories")
	return nil
}

// processEvents fetches and processes events
func (p *Pipeline) processEvents(ctx context.Context) (int, error) {
	p.logger.Info("Processing events")

	// Fetch events from NASA EONET API
	opts := api.FetchEventsOptions{
		Days:   30, // Last 30 days
		Limit:  p.config.ETL.BatchSize,
		Status: "all",
	}

	events, err := p.eonetClient.FetchEvents(ctx, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch events: %w", err)
	}

	// Transform events to database records
	eventRecords := make([]*models.EventRecord, 0, len(events.Events))
	for _, event := range events.Events {
		record, err := p.transformEvent(event)
		if err != nil {
			p.logger.WithError(err).WithField("event_id", event.ID).Warn("Failed to transform event, skipping")
			continue
		}
		eventRecords = append(eventRecords, record)
	}

	// Batch insert events
	if err := p.db.BatchInsertEvents(ctx, eventRecords); err != nil {
		return 0, fmt.Errorf("failed to insert events: %w", err)
	}

	p.logger.WithField("count", len(eventRecords)).Info("Successfully processed events")
	return len(eventRecords), nil
}

// transformEvent transforms an EONET event to a database record
func (p *Pipeline) transformEvent(event models.Event) (*models.EventRecord, error) {
	// Serialize categories to JSON
	categoriesJSON, err := json.Marshal(event.Categories)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal categories: %w", err)
	}

	// Serialize sources to JSON
	sourcesJSON, err := json.Marshal(event.Sources)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sources: %w", err)
	}

	// Serialize geometry to JSON
	geometryJSON, err := json.Marshal(event.Geometry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal geometry: %w", err)
	}

	record := &models.EventRecord{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Link:        event.Link,
		Categories:  string(categoriesJSON),
		Sources:     string(sourcesJSON),
		Geometry:    string(geometryJSON),
		Closed:      event.Closed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return record, nil
}

// HealthCheck performs health checks on all components
func (p *Pipeline) HealthCheck(ctx context.Context) error {
	// Check NASA EONET API
	if err := p.eonetClient.HealthCheck(ctx); err != nil {
		return fmt.Errorf("NASA EONET API health check failed: %w", err)
	}

	// Check database
	if err := p.db.HealthCheck(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetLastRunInfo returns information about the last ETL run
func (p *Pipeline) GetLastRunInfo(ctx context.Context) (*database.ETLRunInfo, error) {
	return p.db.GetLastETLRun(ctx)
}

// Close closes all connections
func (p *Pipeline) Close() error {
	return p.db.Close()
}
