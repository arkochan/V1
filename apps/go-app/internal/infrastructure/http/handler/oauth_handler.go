package handler

import (
	"net/http"
	"user-review-ingest/internal/application/dto"
	"user-review-ingest/internal/application/interfaces"

	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	usecase interfaces.OAuthUsecase
}

func NewOAuthHandler(usecase interfaces.OAuthUsecase) *OAuthHandler {
	return &OAuthHandler{
		usecase: usecase,
	}
}

// @Summary Initiate OAuth Login
// @Description Get the login URL for a specific OAuth provider.
// @Tags OAuth
// @Produce  json
// @Param   provider     path    string  true  "OAuth Provider (e.g., google)"
// @Param   redirect_uri query   string  false "Optional override for the redirect URI"
// @Success 200 {object} dto.OAuthLoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /oauth/{provider}/login [get]
func (h *OAuthHandler) OAuthLogin(c *gin.Context) {
	// Extract provider from URL path
	provider := c.Param("provider")
	redirectURI := c.Query("redirect_uri")

	// Generate a state parameter to prevent CSRF
	state := generateState() // You'll need to implement this function

	loginURL, err := h.usecase.GetLoginURL(c.Request.Context(), provider, redirectURI, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the login URL in JSON response
	response := dto.OAuthLoginResponse{
		URL: loginURL,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary OAuth Callback
// @Description Handle the callback from the OAuth provider after user authorization.
// @Tags OAuth
// @Produce  json
// @Param   provider     path    string  true  "OAuth Provider (e.g., google)"
// @Param   code         query   string  true  "Authorization Code"
// @Param   state        query   string  true  "State"
// @Success 200 {object} dto.OAuthCallbackResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /oauth/{provider}/callback [get]
func (h *OAuthHandler) OAuthCallback(c *gin.Context) {
	// Extract provider from URL path
	provider := c.Param("provider")

	// Parse query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	// Handle the OAuth callback
	user, accessToken, expiresIn, err := h.usecase.HandleCallback(c.Request.Context(), provider, code, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create callback response
	response := dto.OAuthCallbackResponse{
		User:        user,
		AccessToken: accessToken,
		ExpiresIn:   int(expiresIn.Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to generate state parameter (implement as needed)
func generateState() string {
	// Implement a secure random state generator
	// This is a simplified example - in production, use a cryptographically secure random generator
	return "random_state_value"
}

