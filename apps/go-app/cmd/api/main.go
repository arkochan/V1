package main

import (
	"fmt"
	"net/http"
	_ "user-review-ingest/docs" // <-- import generated docs package
	"user-review-ingest/internal/infrastructure/config"
	"user-review-ingest/internal/infrastructure/database"
	"user-review-ingest/internal/infrastructure/http/router"
	"user-review-ingest/internal/infrastructure/observability"
)

// @title User Review Ingest API
// @version 1.0
// @description This is a sample server for a user review ingestion service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Initialize logger
	logger := observability.NewLogger()
	logger.Info().Msg("Starting user-review-ingest service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Setup router
	r := router.SetupRouter(db, logger, cfg)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Info().Msgf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
