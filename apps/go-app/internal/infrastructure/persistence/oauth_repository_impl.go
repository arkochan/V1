package persistence

import (
	"context"
	"time"
	"user-review-ingest/internal/domain/entity"
	"user-review-ingest/internal/domain/repository"
	"user-review-ingest/internal/infrastructure/persistence/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OAuthRepositoryImpl struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewOAuthRepositoryImpl(db *pgxpool.Pool) repository.OAuthRepository {
	return &OAuthRepositoryImpl{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *OAuthRepositoryImpl) FindUserAuthByEmail(ctx context.Context, email string) (*entity.UserAuth, error) {
	user, err := r.queries.GetAuthUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return nil, err
	}

	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return &entity.UserAuth{
		ID:           user.ID.String(),
		Email:        user.Email.String,
		PasswordHash: user.PasswordHash.String,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    deletedAt,
	}, nil
}

func (r *OAuthRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.UserAuth, error) {
	// Convert string ID to UUID
	var userUUID pgtype.UUID
	err := userUUID.Scan(id)
	if err != nil {
		return nil, err
	}

	user, err := r.queries.GetAuthUserByID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return &entity.UserAuth{
		ID:           user.ID.String(),
		Email:        user.Email.String,
		PasswordHash: user.PasswordHash.String,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		DeletedAt:    deletedAt,
	}, nil
}

func (r *OAuthRepositoryImpl) FindOAuthConnectionByProviderID(ctx context.Context, provider, providerID string) (*entity.OAuthConnection, error) {
	oauthConnection, err := r.queries.GetOAuthProviderByProviderID(ctx, sqlc.GetOAuthProviderByProviderIDParams{
		ProviderName:   provider,
		ProviderUserID: providerID,
	})
	if err != nil {
		return nil, err
	}

	return &entity.OAuthConnection{
		ID:             oauthConnection.ID.String(),
		UserID:         oauthConnection.UserID.String(),
		ProviderName:   oauthConnection.ProviderName,
		ProviderUserID: oauthConnection.ProviderUserID,
		Email:          oauthConnection.Email.String,
		Name:           oauthConnection.Name.String,
		AccessToken:    oauthConnection.AccessToken.String,
		RefreshToken:   oauthConnection.RefreshToken.String,
		ExpiresAt:      oauthConnection.ExpiresAt.Time,
		CreatedAt:      oauthConnection.CreatedAt.Time,
		UpdatedAt:      oauthConnection.UpdatedAt.Time,
	}, nil
}

func (r *OAuthRepositoryImpl) CreateUserAuth(ctx context.Context, userAuth *entity.UserAuth) error {
	// Convert string ID to UUID
	var idUUID pgtype.UUID
	err := idUUID.Scan(userAuth.ID)
	if err != nil {
		// If ID is not set, generate a new one
		idUUID = pgtype.UUID{Bytes: uuid.New(), Valid: true}
	}

	_, err = r.queries.CreateAuthUser(ctx, sqlc.CreateAuthUserParams{
		ID:           idUUID,
		Email:        pgtype.Text{String: userAuth.Email, Valid: userAuth.Email != ""},
		PasswordHash: pgtype.Text{String: userAuth.PasswordHash, Valid: userAuth.PasswordHash != ""},
		Status:       userAuth.Status,
	})
	return err
}

func (r *OAuthRepositoryImpl) UpdateUserAuth(ctx context.Context, userAuth *entity.UserAuth) error {
	// Convert string ID to UUID
	var idUUID pgtype.UUID
	err := idUUID.Scan(userAuth.ID)
	if err != nil {
		return err
	}

	_, err = r.queries.UpdateAuthUser(ctx, sqlc.UpdateAuthUserParams{
		ID:        idUUID,
		Email:     pgtype.Text{String: userAuth.Email, Valid: userAuth.Email != ""},
		Status:    userAuth.Status,
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	return err
}

func (r *OAuthRepositoryImpl) CreateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	// Convert string ID to UUID
	var idUUID pgtype.UUID
	err := idUUID.Scan(profile.ID)
	if err != nil {
		// If ID is not set, generate a new one
		idUUID = pgtype.UUID{Bytes: uuid.New(), Valid: true}
	}

	_, err = r.queries.CreateUserProfile(ctx, sqlc.CreateUserProfileParams{
		ID:        idUUID,
		Email:     profile.Email, // Use string directly instead of pgtype.Text
		Name:      pgtype.Text{String: profile.Name, Valid: profile.Name != ""},
		CreatedAt: pgtype.Timestamptz{Time: profile.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: profile.UpdatedAt, Valid: true},
	})
	return err
}

func (r *OAuthRepositoryImpl) UpdateUserProfile(ctx context.Context, profile *entity.UserProfile) error {
	// Convert string ID to UUID
	var idUUID pgtype.UUID
	err := idUUID.Scan(profile.ID)
	if err != nil {
		return err
	}

	_, err = r.queries.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
		ID:        idUUID,
		Email:     profile.Email, // Use string directly instead of pgtype.Text
		Name:      pgtype.Text{String: profile.Name, Valid: profile.Name != ""},
		UpdatedAt: pgtype.Timestamptz{Time: profile.UpdatedAt, Valid: true},
	})
	return err
}

func (r *OAuthRepositoryImpl) FindUserProfileByEmail(ctx context.Context, email string) (*entity.UserProfile, error) {
	profile, err := r.queries.GetUserProfileByEmail(ctx, email) // Use string directly instead of pgtype.Text
	if err != nil {
		return nil, err
	}

	return &entity.UserProfile{
		ID:        profile.ID.String(),
		Email:     profile.Email, // Use string directly
		Name:      profile.Name.String,
		CreatedAt: profile.CreatedAt.Time,
		UpdatedAt: profile.UpdatedAt.Time,
	}, nil
}

func (r *OAuthRepositoryImpl) CreateOAuthConnection(ctx context.Context, connection *entity.OAuthConnection) error {
	// Convert string IDs to UUIDs
	var userIDUUID, idUUID pgtype.UUID

	err := userIDUUID.Scan(connection.UserID)
	if err != nil {
		return err
	}

	err = idUUID.Scan(connection.ID)
	if err != nil {
		// If ID is not set, generate a new one
		idUUID = pgtype.UUID{Bytes: uuid.New(), Valid: true}
	}

	_, err = r.queries.CreateOAuthProvider(ctx, sqlc.CreateOAuthProviderParams{
		ID:             idUUID,
		UserID:         userIDUUID,
		ProviderName:   connection.ProviderName,
		ProviderUserID: connection.ProviderUserID,
		Email:          pgtype.Text{String: connection.Email, Valid: connection.Email != ""},
		Name:           pgtype.Text{String: connection.Name, Valid: connection.Name != ""},
		AccessToken:    pgtype.Text{String: connection.AccessToken, Valid: connection.AccessToken != ""},
		RefreshToken:   pgtype.Text{String: connection.RefreshToken, Valid: connection.RefreshToken != ""},
		ExpiresAt:      pgtype.Timestamptz{Time: connection.ExpiresAt, Valid: true},
		Scopes:         nil, // Scopes can be passed if needed
	})
	return err
}

func (r *OAuthRepositoryImpl) UpdateOAuthConnection(ctx context.Context, connection *entity.OAuthConnection) error {
	// Convert string ID to UUID
	var idUUID pgtype.UUID
	err := idUUID.Scan(connection.ID)
	if err != nil {
		return err
	}

	_, err = r.queries.UpdateOAuthProvider(ctx, sqlc.UpdateOAuthProviderParams{
		ID:           idUUID,
		AccessToken:  pgtype.Text{String: connection.AccessToken, Valid: connection.AccessToken != ""},
		RefreshToken: pgtype.Text{String: connection.RefreshToken, Valid: connection.RefreshToken != ""},
		ExpiresAt:    pgtype.Timestamptz{Time: connection.ExpiresAt, Valid: true},
	})
	return err
}

