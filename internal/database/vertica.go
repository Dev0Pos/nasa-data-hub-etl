package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nasa-data-hub-etl/internal/config"
	"nasa-data-hub-etl/pkg/models"

	"github.com/sirupsen/logrus"
	_ "github.com/vertica/vertica-sql-go"
)

// VerticaDB handles database operations for VerticaDB
type VerticaDB struct {
	db     *sql.DB
	config *config.DatabaseConfig
	logger *logrus.Logger
}

// NewVerticaDB creates a new VerticaDB connection
func NewVerticaDB(cfg *config.DatabaseConfig, logger *logrus.Logger) (*VerticaDB, error) {
	dsn := fmt.Sprintf("vertica://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)

	db, err := sql.Open("vertica", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	verticaDB := &VerticaDB{
		db:     db,
		config: cfg,
		logger: logger,
	}

	// Initialize database schema
	if err := verticaDB.InitializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return verticaDB, nil
}

// InitializeSchema creates the necessary tables if they don't exist
func (v *VerticaDB) InitializeSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			link VARCHAR(500),
			description VARCHAR(10000),
			layers VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id VARCHAR(50) PRIMARY KEY,
			title VARCHAR(500) NOT NULL,
			description VARCHAR(10000),
			link VARCHAR(500),
			categories VARCHAR(10000),
			sources VARCHAR(10000),
			geometry VARCHAR(10000),
			closed VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS etl_runs (
			id BIGINT PRIMARY KEY,
			started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP,
			status VARCHAR(20) NOT NULL,
			events_processed INTEGER DEFAULT 0,
			categories_processed INTEGER DEFAULT 0,
			error_message VARCHAR(10000)
		)`,
	}

	for _, query := range queries {
		if _, err := v.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute schema query: %w", err)
		}
	}

	v.logger.Info("Database schema initialized successfully")
	return nil
}

// InsertEvent inserts or updates an event record
func (v *VerticaDB) InsertEvent(ctx context.Context, event *models.EventRecord) error {
	query := `
		INSERT INTO events (id, title, description, link, categories, sources, geometry, closed, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			link = EXCLUDED.link,
			categories = EXCLUDED.categories,
			sources = EXCLUDED.sources,
			geometry = EXCLUDED.geometry,
			closed = EXCLUDED.closed,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := v.db.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.Description,
		event.Link,
		event.Categories,
		event.Sources,
		event.Geometry,
		event.Closed,
	)

	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	return nil
}

// InsertCategory inserts or updates a category record
func (v *VerticaDB) InsertCategory(ctx context.Context, category *models.CategoryRecord) error {
	query := `
		INSERT INTO categories (id, title, link, description, layers, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			link = EXCLUDED.link,
			description = EXCLUDED.description,
			layers = EXCLUDED.layers,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := v.db.ExecContext(ctx, query,
		category.ID,
		category.Title,
		category.Link,
		category.Description,
		category.Layers,
	)

	if err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}

	return nil
}

// BatchInsertEvents inserts multiple events in a batch
func (v *VerticaDB) BatchInsertEvents(ctx context.Context, events []*models.EventRecord) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			v.logger.WithError(err).Error("Failed to rollback transaction")
		}
	}()

	query := `
		INSERT INTO events (id, title, description, link, categories, sources, geometry, closed, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			link = EXCLUDED.link,
			categories = EXCLUDED.categories,
			sources = EXCLUDED.sources,
			geometry = EXCLUDED.geometry,
			closed = EXCLUDED.closed,
			updated_at = CURRENT_TIMESTAMP
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err := stmt.ExecContext(ctx,
			event.ID,
			event.Title,
			event.Description,
			event.Link,
			event.Categories,
			event.Sources,
			event.Geometry,
			event.Closed,
		)
		if err != nil {
			return fmt.Errorf("failed to execute batch insert: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	v.logger.WithField("count", len(events)).Info("Successfully batch inserted events")
	return nil
}

// BatchInsertCategories inserts multiple categories in a batch
func (v *VerticaDB) BatchInsertCategories(ctx context.Context, categories []*models.CategoryRecord) error {
	if len(categories) == 0 {
		return nil
	}

	tx, err := v.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			v.logger.WithError(err).Error("Failed to rollback transaction")
		}
	}()

	query := `
		INSERT INTO categories (id, title, link, description, layers, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			link = EXCLUDED.link,
			description = EXCLUDED.description,
			layers = EXCLUDED.layers,
			updated_at = CURRENT_TIMESTAMP
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, category := range categories {
		_, err := stmt.ExecContext(ctx,
			category.ID,
			category.Title,
			category.Link,
			category.Description,
			category.Layers,
		)
		if err != nil {
			return fmt.Errorf("failed to execute batch insert: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	v.logger.WithField("count", len(categories)).Info("Successfully batch inserted categories")
	return nil
}

// StartETLRun records the start of an ETL run
func (v *VerticaDB) StartETLRun(ctx context.Context) (int64, error) {
	// VerticaDB doesn't support RETURNING clause, so we'll use a different approach
	// For now, we'll return a timestamp-based ID
	timestamp := time.Now().Unix()

	query := `INSERT INTO etl_runs (id, status) VALUES (?, 'running')`
	_, err := v.db.ExecContext(ctx, query, timestamp)
	if err != nil {
		return 0, fmt.Errorf("failed to start ETL run: %w", err)
	}

	return timestamp, nil
}

// CompleteETLRun records the completion of an ETL run
func (v *VerticaDB) CompleteETLRun(ctx context.Context, runID int64, status string, eventsProcessed, categoriesProcessed int, errorMsg *string) error {
	query := `
		UPDATE etl_runs 
		SET completed_at = CURRENT_TIMESTAMP,
			status = $2,
			events_processed = $3,
			categories_processed = $4,
			error_message = $5
		WHERE id = $1
	`

	_, err := v.db.ExecContext(ctx, query, runID, status, eventsProcessed, categoriesProcessed, errorMsg)
	if err != nil {
		return fmt.Errorf("failed to complete ETL run: %w", err)
	}

	return nil
}

// GetLastETLRun returns information about the last ETL run
func (v *VerticaDB) GetLastETLRun(ctx context.Context) (*ETLRunInfo, error) {
	query := `
		SELECT id, started_at, completed_at, status, events_processed, categories_processed, error_message
		FROM etl_runs
		ORDER BY started_at DESC
		LIMIT 1
	`

	var run ETLRunInfo
	var completedAt sql.NullTime
	var errorMsg sql.NullString

	err := v.db.QueryRowContext(ctx, query).Scan(
		&run.ID,
		&run.StartedAt,
		&completedAt,
		&run.Status,
		&run.EventsProcessed,
		&run.CategoriesProcessed,
		&errorMsg,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No previous runs
		}
		return nil, fmt.Errorf("failed to get last ETL run: %w", err)
	}

	if completedAt.Valid {
		run.CompletedAt = &completedAt.Time
	}
	if errorMsg.Valid {
		run.ErrorMessage = &errorMsg.String
	}

	return &run, nil
}

// ETLRunInfo represents information about an ETL run
type ETLRunInfo struct {
	ID                  int64      `json:"id"`
	StartedAt           time.Time  `json:"started_at"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	Status              string     `json:"status"`
	EventsProcessed     int        `json:"events_processed"`
	CategoriesProcessed int        `json:"categories_processed"`
	ErrorMessage        *string    `json:"error_message,omitempty"`
}

// HealthCheck checks if the database is accessible
func (v *VerticaDB) HealthCheck(ctx context.Context) error {
	return v.db.PingContext(ctx)
}

// Close closes the database connection
func (v *VerticaDB) Close() error {
	return v.db.Close()
}
