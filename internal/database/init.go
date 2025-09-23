package database

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// InitMode represents database initialization mode
type InitMode string

const (
	InitModeCreate InitMode = "Create" // Create database structure
	InitModeRevive InitMode = "Revive" // Skip structure creation
	InitModeAuto   InitMode = "Auto"   // Auto-detect if structure exists
)

// InitializeDatabase creates database structure based on mode
func (v *VerticaDB) InitializeDatabase(ctx context.Context, mode InitMode) error {
	switch mode {
	case InitModeCreate:
		return v.createDatabaseStructure(ctx)
	case InitModeRevive:
		v.logger.Info("Database initialization skipped (revive mode)")
		return nil
	case InitModeAuto:
		return v.autoInitializeDatabase(ctx)
	default:
		return fmt.Errorf("unknown initialization mode: %s, supported modes: Create, Revive, Auto", mode)
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
		description VARCHAR(10000),
		link VARCHAR(1000),
		categories VARCHAR(10000),
		sources VARCHAR(10000),
		geometries VARCHAR(10000),
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
		description VARCHAR(10000),
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

// autoInitializeDatabase automatically detects if database structure exists and creates it if needed
func (v *VerticaDB) autoInitializeDatabase(ctx context.Context) error {
	v.logger.Info("Auto-detecting database structure...")
	
	// Check if events table exists using VerticaDB-specific query
	var tableCount int
	query := `SELECT COUNT(*) FROM v_catalog.tables WHERE table_name = 'events' AND table_schema = 'public'`
	err := v.db.QueryRowContext(ctx, query).Scan(&tableCount)
	if err != nil {
		v.logger.Info("Database structure not found, creating...")
		return v.createDatabaseStructure(ctx)
	}
	
	if tableCount > 0 {
		v.logger.Info("Database structure already exists, skipping creation")
		return nil
	}
	
	v.logger.Info("Database structure not found, creating...")
	return v.createDatabaseStructure(ctx)
}

// ValidateInitMode validates the initialization mode
func ValidateInitMode(mode string) (InitMode, error) {
	normalizedMode := cases.Title(language.English).String(strings.ToLower(mode))
	switch normalizedMode {
	case "Create":
		return InitModeCreate, nil
	case "Revive":
		return InitModeRevive, nil
	case "Auto":
		return InitModeAuto, nil
	default:
		return "", fmt.Errorf("invalid initialization mode: %s, supported modes: Create, Revive, Auto", mode)
	}
}
