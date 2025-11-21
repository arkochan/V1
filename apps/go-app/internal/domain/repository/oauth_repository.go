package repository

import (
	"context"
	"user-review-ingest/internal/domain/entity"
)

type OAuthRepository interface {
	FindByID(ctx context.Context, id string) (*entity.UserAuth, error)
	FindUserAuthByEmail(ctx context.Context, email string) (*entity.UserAuth, error)
	FindOAuthConnectionByProviderID(ctx context.Context, provider, providerID string) (*entity.OAuthConnection, error)
	CreateUserAuth(ctx context.Context, userAuth *entity.UserAuth) error
	UpdateUserAuth(ctx context.Context, userAuth *entity.UserAuth) error
	CreateUserProfile(ctx context.Context, profile *entity.UserProfile) error
	UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error
	FindUserProfileByEmail(ctx context.Context, email string) (*entity.UserProfile, error)
	CreateOAuthConnection(ctx context.Context, connection *entity.OAuthConnection) error
	UpdateOAuthConnection(ctx context.Context, connection *entity.OAuthConnection) error
}
