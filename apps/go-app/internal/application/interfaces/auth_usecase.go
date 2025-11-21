package interfaces

import (
	"context"
	"time"
	"user-review-ingest/internal/domain/entity"
)

type OAuthUsecase interface {
	// OAuth login flow
	GetLoginURL(ctx context.Context, provider, redirectURL, state string) (string, error)

	// OAuth callback handling
	HandleCallback(ctx context.Context, provider, code, state string) (*entity.UserAuth, string, time.Duration, error)

	// Token refresh
	RefreshToken(ctx context.Context, refreshToken, provider string) (string, time.Duration, error)

	// Get user info from provider
	GetUserInfo(ctx context.Context, provider, accessToken string) (*entity.OAuthUser, error)
}

