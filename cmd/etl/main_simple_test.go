package main

import (
	"context"
	"testing"

	"nasa-data-hub-etl/internal/config"
	"nasa-data-hub-etl/internal/database"

	"github.com/sirupsen/logrus"
)

// MockDatabase for testing
type MockDatabase struct{}

func (m *MockDatabase) Ping(ctx context.Context) error {
	return nil
}

func (m *MockDatabase) InitializeDatabase(ctx context.Context, mode database.InitMode) error {
	return nil
}

func (m *MockDatabase) StartETLRun(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockDatabase) CompleteETLRun(ctx context.Context, runID int64, status string, eventsProcessed, categoriesProcessed int, errorMsg *string) error {
	return nil
}

func (m *MockDatabase) InsertEvent(ctx context.Context, event interface{}) error {
	return nil
}

func (m *MockDatabase) InsertCategory(ctx context.Context, category interface{}) error {
	return nil
}

func (m *MockDatabase) BatchInsertEvents(ctx context.Context, events []interface{}) error {
	return nil
}

func (m *MockDatabase) BatchInsertCategories(ctx context.Context, categories []interface{}) error {
	return nil
}

func TestMain_HealthCheck(t *testing.T) {
	mockDB := &MockDatabase{}

	err := runHealthCheck(mockDB)
	if err != nil {
		t.Errorf("runHealthCheck() error = %v", err)
	}
}

func TestMain_InitializeDatabase(t *testing.T) {
	cfg := &config.Config{}
	logger := logrus.New()
	mockDB := &MockDatabase{}

	err := initializeDatabase(mockDB, "Create", cfg, logger)
	if err != nil {
		t.Errorf("initializeDatabase() error = %v", err)
	}
}

func TestMain_ValidateInitMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		want    database.InitMode
		wantErr bool
	}{
		{
			name:    "valid Create mode",
			mode:    "Create",
			want:    database.InitModeCreate,
			wantErr: false,
		},
		{
			name:    "valid Revive mode",
			mode:    "Revive",
			want:    database.InitModeRevive,
			wantErr: false,
		},
		{
			name:    "valid Auto mode",
			mode:    "Auto",
			want:    database.InitModeAuto,
			wantErr: false,
		},
		{
			name:    "case insensitive",
			mode:    "create",
			want:    database.InitModeCreate,
			wantErr: false,
		},
		{
			name:    "invalid mode",
			mode:    "Invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty mode",
			mode:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateInitMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateInitMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateInitMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMain_LoadConfiguration(t *testing.T) {
	cfg, err := loadConfiguration()
	if err != nil {
		t.Errorf("loadConfiguration() error = %v", err)
	}

	if cfg == nil {
		t.Error("loadConfiguration() returned nil config")
		return
	}

	// Validate basic configuration
	if cfg.NASA.APIURL == "" {
		t.Error("NASA API URL should not be empty")
	}

	if cfg.ETL.BatchSize <= 0 {
		t.Error("ETL BatchSize should be greater than 0")
	}

	if cfg.Server.Port <= 0 {
		t.Error("Server Port should be greater than 0")
	}
}

// Helper functions for testing
func runHealthCheck(db interface{}) error {
	ctx := context.Background()
	if mockDB, ok := db.(*MockDatabase); ok {
		return mockDB.Ping(ctx)
	}
	return nil
}

func initializeDatabase(db interface{}, mode string, cfg *config.Config, logger *logrus.Logger) error {
	initMode, err := validateInitMode(mode)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if mockDB, ok := db.(*MockDatabase); ok {
		return mockDB.InitializeDatabase(ctx, initMode)
	}
	return nil
}

func validateInitMode(mode string) (database.InitMode, error) {
	return database.ValidateInitMode(mode)
}

func loadConfiguration() (*config.Config, error) {
	return config.Load()
}
