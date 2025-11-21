package modules

import (
	"os"
	"user-review-ingest/internal/application/usecase"
	"user-review-ingest/internal/infrastructure/http/handler"
	"user-review-ingest/internal/infrastructure/oauth"
	"user-review-ingest/internal/infrastructure/oauth/google"
	"user-review-ingest/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// RegisterOAuthModule sets up the dependencies for the OAuth module and registers its routes.
func RegisterOAuthModule(router *gin.Engine, db *pgxpool.Pool, logger *zerolog.Logger) {
	// Dependencies for OAuth module
	oauthRepo := persistence.NewOAuthRepositoryImpl(db)

	// Provider Registry and Google Provider
	providerRegistry := oauth.NewProviderRegistry()

	// Register Google provider if environment variables are set
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if googleClientID != "" && googleClientSecret != "" && googleRedirectURL != "" {
		googleProvider := google.NewGoogleProvider(googleClientID, googleClientSecret, googleRedirectURL, nil)
		providerRegistry.Register("google", googleProvider)
	}

	oauthUseCase := usecase.NewOAuthUsecase(oauthRepo, providerRegistry, logger)
	oauthHandler := handler.NewOAuthHandler(oauthUseCase)

	// OAuth routes
	router.GET("/oauth/:provider/login", oauthHandler.OAuthLogin)
	router.GET("/oauth/:provider/callback", oauthHandler.OAuthCallback)
}
