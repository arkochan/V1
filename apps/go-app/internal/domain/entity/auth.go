package entity

import (
	"time"
)

type UserAuth struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Don't expose password hash in JSON
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type OAuthConnection struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ProviderName   string    `json:"provider_name"`
	ProviderUserID string    `json:"provider_user_id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	AccessToken    string    `json:"access_token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
	ExpiresAt      time.Time `json:"expires_at,omitempty"`
	Scopes         []string  `json:"scopes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type OAuthUser struct {
	ID    string
	Email string
	Name  string
}
