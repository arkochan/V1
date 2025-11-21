package oauth

import "user-review-ingest/internal/domain/entity"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Provider interface {
	GetAuthURL(redirectURL, state string) (string, error)
	ExchangeToken(code string) (*Token, error)
	RefreshToken(refreshToken string) (*Token, error)
	GetUserInfo(accessToken string) (*entity.OAuthUser, error)
}