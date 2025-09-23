package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewLogger(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Error("New() returned nil")
		return
	}

	// Test that logger is properly configured
	if logger.Formatter == nil {
		t.Error("New() did not set formatter")
	}

	// Test that logger can log messages
	logger.Info("Test log message")
}

func TestNewLogger_WithFields(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Error("NewLogger() returned nil")
		return
	}

	// Test logging with fields
	entry := logger.WithFields(logrus.Fields{
		"test_field": "test_value",
		"number":     42,
	})

	if entry == nil {
		t.Error("WithFields() returned nil")
	}

	entry.Info("Test message with fields")
}

func TestNewLogger_WithField(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Error("NewLogger() returned nil")
		return
	}

	// Test logging with single field
	entry := logger.WithField("test_field", "test_value")

	if entry == nil {
		t.Error("WithField() returned nil")
	}

	entry.Info("Test message with single field")
}

func TestNewLogger_AllLevels(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Error("NewLogger() returned nil")
		return
	}

	// Test all log levels
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
	// Note: Fatal and Panic would exit the program, so we skip them in tests
}

func TestNewLogger_JSONFormatter(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Error("NewLogger() returned nil")
		return
	}

	// Test that formatter is JSON
	if _, ok := logger.Formatter.(*logrus.JSONFormatter); !ok {
		t.Error("NewLogger() should use JSONFormatter")
	}
}
