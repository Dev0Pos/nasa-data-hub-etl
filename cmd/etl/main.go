package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nasa-data-hub-etl/internal/config"
	"nasa-data-hub-etl/internal/database"
	"nasa-data-hub-etl/internal/etl"
	"nasa-data-hub-etl/internal/logger"
	"nasa-data-hub-etl/internal/server"
)

func main() {
	// Parse command line flags
	var (
		healthCheck = flag.Bool("health", false, "Run health check and exit")
		dbInitMode  = flag.String("db-init", "Auto", "Database initialization mode: Create, Revive, or Auto")
	)
	flag.Parse()

	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Create ETL pipeline
	pipeline, err := etl.NewPipeline(cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to create ETL pipeline")
	}
	defer pipeline.Close()

	// Initialize database structure based on mode
	initMode, err := database.ValidateInitMode(*dbInitMode)
	if err != nil {
		log.WithError(err).Fatal("Invalid database initialization mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := pipeline.InitializeDatabase(ctx, initMode); err != nil {
		log.WithError(err).Fatal("Failed to initialize database structure")
	}

	// Handle health check flag
	if *healthCheck {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := pipeline.HealthCheck(ctx); err != nil {
			log.WithError(err).Error("Health check failed")
			os.Exit(1)
		}

		log.Info("Health check passed")
		os.Exit(0)
	}

	// Setup graceful shutdown
	cancel()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// Start HTTP server for health checks and metrics
	httpServer := server.NewServer(cfg, pipeline, log)
	go func() {
		if err := httpServer.Start(ctx); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server failed")
		}
	}()

	// Start ETL pipeline
	log.Info("Starting NASA Data Hub ETL pipeline")
	if err := pipeline.Run(ctx); err != nil {
		log.WithError(err).Fatal("ETL pipeline failed")
	}

	log.Info("ETL pipeline completed successfully")
}
