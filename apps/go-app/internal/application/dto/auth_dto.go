package dto

type OAuthLoginRequest struct {
	Provider string `json:"provider" validate:"required"`
	Redirect string `json:"redirect,omitempty"`
	State    string `json:"state,omitempty"`
}

type OAuthLoginResponse struct {
	URL string `json:"url"`
}

type OAuthCallbackRequest struct {
	Provider string `json:"provider" validate:"required"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state,omitempty"`
	Redirect string `json:"redirect,omitempty"`
}

type OAuthCallbackResponse struct {
	User        interface{} `json:"user"`
	AccessToken string      `json:"access_token"`
	ExpiresIn   int         `json:"expires_in"`
	RedirectURL string      `json:"redirect_url,omitempty"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
}

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

