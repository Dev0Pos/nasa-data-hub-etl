package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"NASA_API_URL",
		"NASA_API_KEY",
		"ETL_BATCH_SIZE",
		"ETL_INTERVAL",
		"ETL_RETRY_ATTEMPTS",
		"ETL_RETRY_DELAY",
		"SERVER_PORT",
		"SERVER_READ_TIMEOUT",
		"SERVER_WRITE_TIMEOUT",
		"LOG_LEVEL",
	}

	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
	}

	// Clean up after test
	defer func() {
		for _, env := range envVars {
			if val, exists := originalEnv[env]; exists {
				os.Setenv(env, val)
			} else {
				os.Unsetenv(env)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "default configuration",
			envVars: map[string]string{},
			wantErr: false,
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"NASA_API_URL":      "https://custom.api.com",
				"ETL_BATCH_SIZE":    "2000",
				"ETL_INTERVAL":      "2h",
				"ETL_RETRY_ATTEMPTS": "5",
				"ETL_RETRY_DELAY":   "60s",
				"SERVER_PORT":       "9090",
				"LOG_LEVEL":         "debug",
			},
			wantErr: false,
		},
		{
			name: "invalid batch size",
			envVars: map[string]string{
				"ETL_BATCH_SIZE": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid retry attempts",
			envVars: map[string]string{
				"ETL_RETRY_ATTEMPTS": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid server port",
			envVars: map[string]string{
				"SERVER_PORT": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up environment after test
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if cfg == nil {
					t.Error("Load() returned nil config")
					return
				}

				// Validate configuration values
				if tt.envVars["NASA_API_URL"] != "" && cfg.NASA.APIURL != tt.envVars["NASA_API_URL"] {
					t.Errorf("Load() NASA.APIURL = %v, want %v", cfg.NASA.APIURL, tt.envVars["NASA_API_URL"])
				}

				if tt.envVars["ETL_BATCH_SIZE"] != "" {
					expectedBatchSize := 2000
					if cfg.ETL.BatchSize != expectedBatchSize {
						t.Errorf("Load() ETL.BatchSize = %v, want %v", cfg.ETL.BatchSize, expectedBatchSize)
					}
				}

				if tt.envVars["SERVER_PORT"] != "" {
					expectedPort := 9090
					if cfg.Server.Port != expectedPort {
						t.Errorf("Load() Server.Port = %v, want %v", cfg.Server.Port, expectedPort)
					}
				}
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				NASA: NASAConfig{
					APIURL: "https://eonet.gsfc.nasa.gov/api/v3",
				},
				ETL: ETLConfig{
					BatchSize:     1000,
					Interval:      time.Hour,
					RetryAttempts: 3,
					RetryDelay:    30 * time.Second,
				},
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "empty NASA API URL",
			config: &Config{
				NASA: NASAConfig{
					APIURL: "",
				},
				ETL: ETLConfig{
					BatchSize:     1000,
					Interval:      time.Hour,
					RetryAttempts: 3,
					RetryDelay:    30 * time.Second,
				},
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid batch size",
			config: &Config{
				NASA: NASAConfig{
					APIURL: "https://eonet.gsfc.nasa.gov/api/v3",
				},
				ETL: ETLConfig{
					BatchSize:     0,
					Interval:      time.Hour,
					RetryAttempts: 3,
					RetryDelay:    30 * time.Second,
				},
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid retry attempts",
			config: &Config{
				NASA: NASAConfig{
					APIURL: "https://eonet.gsfc.nasa.gov/api/v3",
				},
				ETL: ETLConfig{
					BatchSize:     1000,
					Interval:      time.Hour,
					RetryAttempts: 0,
					RetryDelay:    30 * time.Second,
				},
				Server: ServerConfig{
					Port:         8080,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid server port",
			config: &Config{
				NASA: NASAConfig{
					APIURL: "https://eonet.gsfc.nasa.gov/api/v3",
				},
				ETL: ETLConfig{
					BatchSize:     1000,
					Interval:      time.Hour,
					RetryAttempts: 3,
					RetryDelay:    30 * time.Second,
				},
				Server: ServerConfig{
					Port:         0,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_LoadSecrets(t *testing.T) {
	// Save original environment
	originalAPIKey := os.Getenv("NASA_API_KEY")

	// Clean up after test
	defer func() {
		if originalAPIKey != "" {
			os.Setenv("NASA_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("NASA_API_KEY")
		}
	}()

	tests := []struct {
		name     string
		apiKey   string
		expected string
	}{
		{
			name:     "no API key",
			apiKey:   "",
			expected: "",
		},
		{
			name:     "with API key",
			apiKey:   "test-api-key",
			expected: "test-api-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.apiKey != "" {
				os.Setenv("NASA_API_KEY", tt.apiKey)
			} else {
				os.Unsetenv("NASA_API_KEY")
			}

			config := &Config{}
			config.LoadSecrets()

			if config.NASA.APIKey != tt.expected {
				t.Errorf("LoadSecrets() APIKey = %v, want %v", config.NASA.APIKey, tt.expected)
			}
		})
	}
}
