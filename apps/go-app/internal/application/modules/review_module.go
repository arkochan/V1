package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"user-review-ingest/internal/application/usecase"
	"user-review-ingest/internal/infrastructure/http/handler"
	"user-review-ingest/internal/infrastructure/persistence"
)

// RegisterReviewModule sets up the dependencies for the review module and registers its routes.
func RegisterReviewModule(router *gin.RouterGroup, db *pgxpool.Pool) {
	// Dependencies for Review module
	reviewRepo := persistence.NewReviewRepositoryImpl(db)
	reviewUseCase := usecase.NewReviewUseCaseImpl(reviewRepo)
	reviewHandler := handler.NewReviewHandler(reviewUseCase)

	// Review routes
	reviews := router.Group("/reviews")
	{
		reviews.POST("", reviewHandler.CreateReview)
		reviews.GET("/:id", reviewHandler.GetReview)
		reviews.PUT("/:id", reviewHandler.UpdateReview)
		reviews.DELETE("/:id", reviewHandler.DeleteReview)
		reviews.GET("", reviewHandler.ListReviews)
	}
}
