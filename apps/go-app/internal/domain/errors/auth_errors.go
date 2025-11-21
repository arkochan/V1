package errors

import "errors"

var (
	// OAuth errors
	ErrInvalidProvider        = errors.New("invalid oauth provider")
	ErrOAuthRequestFailed     = errors.New("oauth request failed")
	ErrTokenExchangeFailed    = errors.New("token exchange failed")
	ErrInvalidToken           = errors.New("invalid token")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserCreationFailed     = errors.New("user creation failed")
	ErrStateMismatch          = errors.New("state parameter mismatch")
	ErrMissingState           = errors.New("missing state parameter")
	ErrMissingCode            = errors.New("missing authorization code")
	ErrInvalidRedirectURL     = errors.New("invalid redirect URL")
)