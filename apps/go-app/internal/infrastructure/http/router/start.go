package router

import (
	"user-review-ingest/internal/application/modules"
	"user-review-ingest/internal/infrastructure/config"
	"user-review-ingest/internal/infrastructure/http/handler"
	"user-review-ingest/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *pgxpool.Pool, logger *zerolog.Logger, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// Health Check and Swagger
	healthHandler := handler.NewHealthHandler()
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	modules.RegisterOAuthModule(r, db, logger)

	// Versioned API Group
	v1RouterGroup := r.Group("/v1")
	{
		modules.RegisterReviewModule(v1RouterGroup, db)
	}

	return r
}
