package usecase

import (
	"context"
	"fmt"
	"time"
	"user-review-ingest/internal/application/interfaces"
	"user-review-ingest/internal/domain/entity"
	"user-review-ingest/internal/domain/repository"
	"user-review-ingest/internal/infrastructure/oauth"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type oauthUsecase struct {
	oauthRepo        repository.OAuthRepository
	providerRegistry *oauth.ProviderRegistry
	logger           *zerolog.Logger
}

func NewOAuthUsecase(oauthRepo repository.OAuthRepository, providerRegistry *oauth.ProviderRegistry, logger *zerolog.Logger) interfaces.OAuthUsecase {
	return &oauthUsecase{
		oauthRepo:        oauthRepo,
		providerRegistry: providerRegistry,
		logger:           logger,
	}
}

func (uc *oauthUsecase) GetLoginURL(ctx context.Context, providerName, redirectURL, state string) (string, error) {
	provider, err := uc.providerRegistry.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	return provider.GetAuthURL(redirectURL, state)
}

func (uc *oauthUsecase) HandleCallback(ctx context.Context, providerName, code, state string) (*entity.UserAuth, string, time.Duration, error) {
	fmt.Printf("DEBUG: Starting HandleCallback - Provider: %s, Code: %s, State: %s\n", providerName, code, state)

	provider, err := uc.providerRegistry.GetProvider(providerName)
	if err != nil {
		fmt.Printf("DEBUG: Error getting provider: %v\n", err)
		return nil, "", 0, err
	}

	token, err := provider.ExchangeToken(code)
	if err != nil {
		fmt.Printf("DEBUG: Error exchanging token: %v\n", err)
		return nil, "", 0, err
	}
	fmt.Printf("DEBUG: Token received - Access: %s, Refresh: %s, ExpiresIn: %d\n",
		token.AccessToken[:min(10, len(token.AccessToken))],
		token.RefreshToken[:min(10, len(token.RefreshToken))],
		token.ExpiresIn)

	oauthUser, err := provider.GetUserInfo(token.AccessToken)
	if err != nil {
		fmt.Printf("DEBUG: Error getting user info: %v\n", err)
		return nil, "", 0, err
	}
	fmt.Printf("DEBUG: OAuth User - ID: %s, Email: %s, Name: %s\n", oauthUser.ID, oauthUser.Email, oauthUser.Name)

	// Check if OAuth connection already exists
	oauthConn, err := uc.oauthRepo.FindOAuthConnectionByProviderID(ctx, providerName, oauthUser.ID)
	if err != nil {
		fmt.Printf("DEBUG: No existing OAuth connection found for provider: %s, user ID: %s\n", providerName, oauthUser.ID)

		// Check if user profile exists by email
		profile, err := uc.oauthRepo.FindUserProfileByEmail(ctx, oauthUser.Email)
		if err != nil {
			fmt.Printf("DEBUG: No existing profile found for email: %s, creating new profile\n", oauthUser.Email)

			// Create new user profile
			profile = &entity.UserProfile{
				ID:        uuid.New().String(),
				Email:     oauthUser.Email,
				Name:      oauthUser.Name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := uc.oauthRepo.CreateUserProfile(ctx, profile); err != nil {
				uc.logger.Error().Err(err).Msg("failed to create user profile")
				return nil, "", 0, fmt.Errorf("failed to create user profile: %w", err)
			}
		} else {
			fmt.Printf("DEBUG: Found existing profile - ID: %s, Email: %s, Name: %s\n", profile.ID, profile.Email, profile.Name)
			// Profile exists, update it
			profile.Name = oauthUser.Name
			profile.UpdatedAt = time.Now()

			if err := uc.oauthRepo.UpdateUserProfile(ctx, profile); err != nil {
				return nil, "", 0, err
			}
		}

		// Check if user auth record already exists by email
		existingUserAuth, err := uc.oauthRepo.FindUserAuthByEmail(ctx, profile.Email)
		if err != nil {
			// Could not find user auth by email. This implies it doesn't exist.
			// Let's create it, using the profile ID we found or created earlier.
			fmt.Printf("DEBUG: No existing user auth found for email: %s, creating new auth with profile ID: %s\n", profile.Email, profile.ID)
			userAuth := &entity.UserAuth{
				ID:           profile.ID,
				Email:        profile.Email,
				PasswordHash: "", // No password for OAuth users
				Status:       "active",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if err := uc.oauthRepo.CreateUserAuth(ctx, userAuth); err != nil {
				uc.logger.Error().Err(err).Msg("failed to create user auth")
				return nil, "", 0, fmt.Errorf("failed to create user auth: %w", err)
			}
			existingUserAuth = userAuth
		} else {
			// Found user auth by email. The profile might need to be linked or updated.
			fmt.Printf("DEBUG: Found existing user auth by email - ID: %s, Email: %s\n", existingUserAuth.ID, existingUserAuth.Email)

			// Update the existing user auth with new information from profile (if any)
			existingUserAuth.Email = profile.Email // Should be the same
			existingUserAuth.UpdatedAt = time.Now()

			if err := uc.oauthRepo.UpdateUserAuth(ctx, existingUserAuth); err != nil {
				return nil, "", 0, err
			}
		}

		// Create new OAuth connection
		oauthConn = &entity.OAuthConnection{
			ID:             uuid.New().String(),
			UserID:         existingUserAuth.ID,
			ProviderName:   providerName,
			ProviderUserID: oauthUser.ID,
			Email:          oauthUser.Email,
			Name:           oauthUser.Name,
			AccessToken:    token.AccessToken,
			RefreshToken:   token.RefreshToken,
			ExpiresAt:      time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := uc.oauthRepo.CreateOAuthConnection(ctx, oauthConn); err != nil {
			uc.logger.Error().Err(err).Msg("failed to create oauth connection")
			return nil, "", 0, fmt.Errorf("failed to create oauth connection: %w", err)
		}

		fmt.Printf("DEBUG: Created new OAuth connection - ID: %s, UserID: %s\n", oauthConn.ID, oauthConn.UserID)
		return existingUserAuth, token.AccessToken, time.Duration(token.ExpiresIn) * time.Second, nil
	} else {
		fmt.Printf("DEBUG: Found existing OAuth connection - ID: %s, UserID: %s, ProviderUserID: %s\n",
			oauthConn.ID, oauthConn.UserID, oauthConn.ProviderUserID)

		// OAuth connection exists, update tokens
		oauthConn.AccessToken = token.AccessToken
		oauthConn.RefreshToken = token.RefreshToken
		oauthConn.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		oauthConn.UpdatedAt = time.Now()

		if err := uc.oauthRepo.UpdateOAuthConnection(ctx, oauthConn); err != nil {
			return nil, "", 0, err
		}

		// Get the associated user auth
		userAuth, err := uc.oauthRepo.FindByID(ctx, oauthConn.UserID)
		if err != nil {
			return nil, "", 0, err
		}

		fmt.Printf("DEBUG: Updated existing OAuth connection and retrieved user auth - ID: %s\n", userAuth.ID)
		return userAuth, token.AccessToken, time.Duration(token.ExpiresIn) * time.Second, nil
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (uc *oauthUsecase) RefreshToken(ctx context.Context, refreshToken, providerName string) (string, time.Duration, error) {
	provider, err := uc.providerRegistry.GetProvider(providerName)
	if err != nil {
		return "", 0, err
	}

	newToken, err := provider.RefreshToken(refreshToken)
	if err != nil {
		return "", 0, err
	}

	return newToken.AccessToken, time.Duration(newToken.ExpiresIn) * time.Second, nil
}

func (uc *oauthUsecase) GetUserInfo(ctx context.Context, provider, accessToken string) (*entity.OAuthUser, error) {
	providerObj, err := uc.providerRegistry.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	return providerObj.GetUserInfo(accessToken)
}
