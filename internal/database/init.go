package database

import (
	"context"
	"fmt"
	"strings"
)

// InitMode represents database initialization mode
type InitMode string

const (
	InitModeCreate InitMode = "Create" // Create database structure
	InitModeRevive InitMode = "Revive" // Skip structure creation
)

// InitializeDatabase creates database structure based on mode
func (v *VerticaDB) InitializeDatabase(ctx context.Context, mode InitMode) error {
	switch mode {
	case InitModeCreate:
		return v.createDatabaseStructure(ctx)
	case InitModeRevive:
		v.logger.Info("Database initialization skipped (revive mode)")
		return nil
	default:
		return fmt.Errorf("unknown initialization mode: %s, supported modes: Create, Revive", mode)
	}
}

// createDatabaseStructure creates all required tables and indexes
func (v *VerticaDB) createDatabaseStructure(ctx context.Context) error {
	v.logger.Info("Creating database structure...")

	// Create events table
	if err := v.createEventsTable(ctx); err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}

	// Create categories table
	if err := v.createCategoriesTable(ctx); err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	// Create indexes
	if err := v.createIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	v.logger.Info("Database structure created successfully")
	return nil
}

// createEventsTable creates the events table
func (v *VerticaDB) createEventsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id VARCHAR(255) PRIMARY KEY,
		title VARCHAR(500),
		description TEXT,
		link VARCHAR(1000),
		categories TEXT,
		sources TEXT,
		geometries TEXT,
		date TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := v.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute events table creation: %w", err)
	}

	v.logger.Info("Events table created successfully")
	return nil
}

// createCategoriesTable creates the categories table
func (v *VerticaDB) createCategoriesTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id INT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := v.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute categories table creation: %w", err)
	}

	v.logger.Info("Categories table created successfully")
	return nil
}

// createIndexes creates performance indexes
func (v *VerticaDB) createIndexes(ctx context.Context) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_events_date ON events(date)",
		"CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_categories_title ON categories(title)",
	}

	for _, indexQuery := range indexes {
		if _, err := v.db.ExecContext(ctx, indexQuery); err != nil {
			return fmt.Errorf("failed to create index '%s': %w", indexQuery, err)
		}
	}

	v.logger.Info("Database indexes created successfully")
	return nil
}

// ValidateInitMode validates the initialization mode
func ValidateInitMode(mode string) (InitMode, error) {
	normalizedMode := strings.Title(strings.ToLower(mode))
	switch normalizedMode {
	case "Create":
		return InitModeCreate, nil
	case "Revive":
		return InitModeRevive, nil
	default:
		return "", fmt.Errorf("invalid initialization mode: %s, supported modes: Create, Revive", mode)
	}
}
