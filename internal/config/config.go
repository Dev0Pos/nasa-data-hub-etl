package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	NASA     NASAConfig     `mapstructure:"nasa"`
	Database DatabaseConfig `mapstructure:"database"`
	ETL      ETLConfig      `mapstructure:"etl"`
	Server   ServerConfig   `mapstructure:"server"`
}

// NASAConfig holds NASA EONET API configuration
type NASAConfig struct {
	APIURL string `mapstructure:"api_url"`
	APIKey string `mapstructure:"api_key"`
}

// DatabaseConfig holds VerticaDB configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// ETLConfig holds ETL pipeline configuration
type ETLConfig struct {
	BatchSize     int           `mapstructure:"batch_size"`
	Interval      time.Duration `mapstructure:"interval"`
	RetryAttempts int           `mapstructure:"retry_attempts"`
	RetryDelay    time.Duration `mapstructure:"retry_delay"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/nasa-data-hub-etl")

	// Set default values
	setDefaults()

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Load secrets from environment variables
	config.LoadSecrets()

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// NASA API defaults
	viper.SetDefault("nasa.api_url", "https://eonet.gsfc.nasa.gov/api/v3")
	viper.SetDefault("nasa.api_key", "")

	// ETL defaults
	viper.SetDefault("etl.batch_size", 1000)
	viper.SetDefault("etl.interval", "1h")
	viper.SetDefault("etl.retry_attempts", 3)
	viper.SetDefault("etl.retry_delay", "30s")

	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
}

// LoadSecrets loads sensitive configuration from environment variables
func (c *Config) LoadSecrets() {
	// Load database configuration from environment
	if host := os.Getenv("DATABASE_HOST"); host != "" {
		c.Database.Host = host
	}
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Database.Port = p
		}
	}
	if database := os.Getenv("DATABASE_NAME"); database != "" {
		c.Database.Database = database
	}
	if username := os.Getenv("DATABASE_USERNAME"); username != "" {
		c.Database.Username = username
	}
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		c.Database.Password = password
	}

	// Load NASA API key from environment
	if apiKey := os.Getenv("NASA_API_KEY"); apiKey != "" {
		c.NASA.APIKey = apiKey
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.NASA.APIURL == "" {
		return fmt.Errorf("nasa.api_url is required")
	}

	// Database configuration is loaded from environment variables
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required (set via DATABASE_HOST environment variable)")
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("database.port must be between 1 and 65535 (set via DATABASE_PORT environment variable)")
	}

	if c.Database.Database == "" {
		return fmt.Errorf("database.database is required (set via DATABASE_NAME environment variable)")
	}

	if c.Database.Username == "" {
		return fmt.Errorf("database.username is required (set via DATABASE_USERNAME environment variable)")
	}

	if c.Database.Password == "" {
		return fmt.Errorf("database.password is required (set via DATABASE_PASSWORD environment variable)")
	}

	if c.ETL.BatchSize <= 0 {
		return fmt.Errorf("etl.batch_size must be greater than 0")
	}

	if c.ETL.RetryAttempts < 0 {
		return fmt.Errorf("etl.retry_attempts must be non-negative")
	}

	return nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("vertica://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.SSLMode,
	)
}
