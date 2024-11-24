package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// KeycloakConfig holds the environment variables required for Keycloak requests.
type KeycloakConfig struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
}

// LoadConfig loads and validates Keycloak configuration from environment variables.
func LoadConfig() (*KeycloakConfig, error) {
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	keycloakRealm := os.Getenv("KEYCLOAK_REALM")
	keycloakClientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	if keycloakURL == "" || keycloakRealm == "" || keycloakClientID == "" || keycloakSecret == "" {
		return nil, fmt.Errorf("missing one or more required environment variables: KEYCLOAK_URL, KEYCLOAK_REALM, KEYCLOAK_CLIENT_ID, KEYCLOAK_CLIENT_SECRET")
	}

	return &KeycloakConfig{
		URL:          keycloakURL,
		Realm:        keycloakRealm,
		ClientID:     keycloakClientID,
		ClientSecret: keycloakSecret,
	}, nil
}

// MakeRequest sends an HTTP request to Keycloak and handles the response.
func MakeRequest(ctx context.Context, method, endpoint string, data url.Values) ([]byte, int, error) {
	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set up the HTTP client
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, nil
}
