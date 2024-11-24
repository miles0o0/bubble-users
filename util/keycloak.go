package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/miles0o0/bubble-users/graph/model"
)

func KeycloakLogin(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	// Load Keycloak configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Build form data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("username", username)
	data.Set("password", password)

	// Make the request
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", config.URL, config.Realm)
	body, statusCode, err := MakeRequest(ctx, "POST", tokenURL, data)
	if err != nil {
		return nil, err
	}

	// Handle non-200 responses
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to login to Keycloak: %s", http.StatusText(statusCode))
	}

	// Parse response
	var tokenResponse model.LoginResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Keycloak response: %w", err)
	}

	return &tokenResponse, nil
}

func KeycloakLogout(ctx context.Context, refreshToken string) (bool, error) {
	// Load Keycloak configuration
	config, err := LoadConfig()
	if err != nil {
		return false, err
	}

	// Build form data
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("refresh_token", refreshToken)

	// Make the request
	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", config.URL, config.Realm)
	_, statusCode, err := MakeRequest(ctx, "POST", logoutURL, data)
	if err != nil {
		return false, err
	}

	// Handle non-204 responses
	if statusCode != http.StatusNoContent {
		return false, fmt.Errorf("failed to logout from Keycloak: %s", http.StatusText(statusCode))
	}

	return true, nil
}

func KeycloakRefresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	// Load Keycloak configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Build form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("refresh_token", refreshToken)

	// Make the request
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", config.URL, config.Realm)
	body, statusCode, err := MakeRequest(ctx, "POST", tokenURL, data)
	if err != nil {
		return nil, err
	}

	// Handle non-200 responses
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh token: %s", http.StatusText(statusCode))
	}

	// Parse response
	var tokenResponse model.LoginResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse Keycloak response: %w", err)
	}

	return &tokenResponse, nil
}
