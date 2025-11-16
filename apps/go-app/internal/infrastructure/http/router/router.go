package router

import (
	_ "user-review-ingest/docs"
	"user-review-ingest/internal/application/usecase"
	"user-review-ingest/internal/infrastructure/http/handler"
	"user-review-ingest/internal/infrastructure/http/middleware"
	"user-review-ingest/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *pgxpool.Pool, logger *zerolog.Logger) *gin.Engine {
	r := gin.New()

	// Middlewares
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	// r.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // Example

	// Repositories
	reviewRepo := persistence.NewReviewRepositoryImpl(db)

	// Use Cases
	reviewUseCase := usecase.NewReviewUseCaseImpl(reviewRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler()
	reviewHandler := handler.NewReviewHandler(
		reviewUseCase,
	)

	// Routes
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	{
		reviews := v1.Group("/reviews")
		{
			reviews.POST("", reviewHandler.CreateReview)
			reviews.GET("/:id", reviewHandler.GetReview)
			reviews.PUT("/:id", reviewHandler.UpdateReview)
			reviews.DELETE("/:id", reviewHandler.DeleteReview)
			reviews.GET("", reviewHandler.ListReviews)
		}
	}

	return r
}
