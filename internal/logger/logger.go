package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New creates a new logger instance
func New() *logrus.Logger {
	log := logrus.New()

	// Set JSON formatter for structured logging
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function",
		},
	})

	// Set output to stdout
	log.SetOutput(os.Stdout)

	// Set log level based on environment
	if os.Getenv("LOG_LEVEL") != "" {
		level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err == nil {
			log.SetLevel(level)
		} else {
			log.SetLevel(logrus.InfoLevel)
		}
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
