package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"nasa-data-hub-etl/internal/config"
	"nasa-data-hub-etl/internal/etl"

	"github.com/sirupsen/logrus"
)

// Server represents the HTTP server
type Server struct {
	config   *config.Config
	pipeline *etl.Pipeline
	logger   *logrus.Logger
	server   *http.Server
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, pipeline *etl.Pipeline, logger *logrus.Logger) *Server {
	return &Server{
		config:   cfg,
		pipeline: pipeline,
		logger:   logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Health check endpoints
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ready", s.readyHandler)
	mux.HandleFunc("/metrics", s.metricsHandler)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:      mux,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	s.logger.WithField("port", s.config.Server.Port).Info("Starting HTTP server")

	go func() {
		<-ctx.Done()
		s.logger.Info("Shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			s.logger.WithError(err).Error("Failed to shutdown HTTP server gracefully")
		}
	}()

	return s.server.ListenAndServe()
}

// healthHandler handles health check requests
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Perform health checks
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := s.pipeline.HealthCheck(ctx); err != nil {
		s.logger.WithError(err).Error("Health check failed")
		http.Error(w, "Health check failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}

// readyHandler handles readiness check requests
func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, readiness is the same as health
	// In the future, this could check if the application is ready to process requests
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.pipeline.HealthCheck(ctx); err != nil {
		s.logger.WithError(err).Error("Readiness check failed")
		http.Error(w, "Not ready", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ready","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}

// metricsHandler handles metrics requests
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, return basic metrics
	// In the future, this could integrate with Prometheus
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// Get last run info
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	lastRun, err := s.pipeline.GetLastRunInfo(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get last run info")
		fmt.Fprintf(w, "# Error getting metrics: %v\n", err)
		return
	}

	if lastRun != nil {
		fmt.Fprintf(w, "# HELP etl_runs_total Total number of ETL runs\n")
		fmt.Fprintf(w, "# TYPE etl_runs_total counter\n")
		fmt.Fprintf(w, "etl_runs_total{status=\"%s\"} 1\n", lastRun.Status)

		fmt.Fprintf(w, "# HELP etl_events_processed_total Total events processed\n")
		fmt.Fprintf(w, "# TYPE etl_events_processed_total counter\n")
		fmt.Fprintf(w, "etl_events_processed_total %d\n", lastRun.EventsProcessed)

		fmt.Fprintf(w, "# HELP etl_categories_processed_total Total categories processed\n")
		fmt.Fprintf(w, "# TYPE etl_categories_processed_total counter\n")
		fmt.Fprintf(w, "etl_categories_processed_total %d\n", lastRun.CategoriesProcessed)
	} else {
		fmt.Fprintf(w, "# No ETL runs found\n")
	}
}
