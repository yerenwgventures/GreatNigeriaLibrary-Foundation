package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthManager manages OAuth authentication
type OAuthManager struct {
	logger Logger
	config Config
}

// Logger interface for OAuth manager
type Logger interface {
	Info(msg string)
	Error(msg string)
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// Config interface for OAuth manager
type Config interface {
	GetOAuthConfig(provider string) *oauth2.Config
}

// OAuthUserInfo represents OAuth user information
type OAuthUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

// GoogleUserInfo represents Google OAuth user info
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager(logger Logger, config Config) *OAuthManager {
	return &OAuthManager{
		logger: logger,
		config: config,
	}
}

// GetAuthURL returns the OAuth authorization URL
func (o *OAuthManager) GetAuthURL(provider, state string) (string, error) {
	config := o.getOAuthConfig(provider)
	if config == nil {
		return "", errors.New("unsupported OAuth provider")
	}

	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	o.logger.WithFields(map[string]interface{}{
		"provider": provider,
		"state":    state,
	}).Info("Generated OAuth authorization URL")

	return url, nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (o *OAuthManager) ExchangeCodeForToken(provider, code string) (*oauth2.Token, error) {
	config := o.getOAuthConfig(provider)
	if config == nil {
		return nil, errors.New("unsupported OAuth provider")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		o.logger.WithError(err).WithField("provider", provider).Error("Failed to exchange code for token")
		return nil, err
	}

	o.logger.WithField("provider", provider).Info("Successfully exchanged code for token")
	return token, nil
}

// GetUserInfo retrieves user information using OAuth token
func (o *OAuthManager) GetUserInfo(provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	switch provider {
	case "google":
		return o.getGoogleUserInfo(token)
	default:
		return nil, errors.New("unsupported OAuth provider")
	}
}

// getGoogleUserInfo retrieves Google user information
func (o *OAuthManager) getGoogleUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		o.logger.WithError(err).Error("Failed to get Google user info")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.logger.WithField("status_code", resp.StatusCode).Error("Google API returned error")
		return nil, fmt.Errorf("Google API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		o.logger.WithError(err).Error("Failed to unmarshal Google user info")
		return nil, err
	}

	userInfo := &OAuthUserInfo{
		ID:       googleUser.ID,
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Picture:  googleUser.Picture,
		Provider: "google",
	}

	o.logger.WithField("email", userInfo.Email).Info("Retrieved Google user info")
	return userInfo, nil
}

// getOAuthConfig returns OAuth configuration for provider
func (o *OAuthManager) getOAuthConfig(provider string) *oauth2.Config {
	switch provider {
	case "google":
		return &oauth2.Config{
			ClientID:     "your-google-client-id",     // Should come from config
			ClientSecret: "your-google-client-secret", // Should come from config
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	default:
		return nil
	}
}

// ValidateState validates OAuth state parameter
func (o *OAuthManager) ValidateState(receivedState, expectedState string) bool {
	return receivedState == expectedState
}

// RevokeToken revokes an OAuth token
func (o *OAuthManager) RevokeToken(provider string, token *oauth2.Token) error {
	switch provider {
	case "google":
		return o.revokeGoogleToken(token)
	default:
		return errors.New("unsupported OAuth provider")
	}
}

// revokeGoogleToken revokes a Google OAuth token
func (o *OAuthManager) revokeGoogleToken(token *oauth2.Token) error {
	revokeURL := fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token.AccessToken)
	
	resp, err := http.Post(revokeURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		o.logger.WithError(err).Error("Failed to revoke Google token")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.logger.WithField("status_code", resp.StatusCode).Error("Google token revocation failed")
		return fmt.Errorf("token revocation failed with status %d", resp.StatusCode)
	}

	o.logger.Info("Successfully revoked Google token")
	return nil
}
