package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/miles0o0/bubble-users/graph/model"
)

func keycloakLogin(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	// **Step 1: Load and validate environment variables**
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	keycloakRealm := os.Getenv("KEYCLOAK_REALM")
	keycloakClientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	// Validate that all required variables are set
	if keycloakURL == "" || keycloakRealm == "" || keycloakClientID == "" || keycloakSecret == "" {
		return nil, fmt.Errorf("missing one or more required environment variables: KEYCLOAK_URL, KEYCLOAK_REALM, KEYCLOAK_CLIENT_ID, KEYCLOAK_CLIENT_SECRET")
	}

	// **Step 2: Construct the token URL**
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, keycloakRealm)

	// **Step 3: Build form data using variables**
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", keycloakClientID)
	data.Set("client_secret", keycloakSecret)
	data.Set("username", username)
	data.Set("password", password)

	// **Step 4: Create the HTTP request**
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// **Step 5: Set up the HTTP client**
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// **Step 6: Execute the request**
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// **Step 7: Read and log the response body**
	body, _ := io.ReadAll(resp.Body)

	// **Step 8: Handle non-200 responses**
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to login to Keycloak: %s", resp.Status)
	}

	// **Step 9: Decode the JSON response**
	var tokenResponse model.LoginResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Printf("Failed to parse Keycloak response. Error: %v", err)
		return nil, fmt.Errorf("failed to parse Keycloak response: %w", err)
	}

	// **Step 10: Return the token response**
	return &tokenResponse, nil
}

func keycloakLogout(ctx context.Context, refreshToken string) (bool, error) {
	// **Step 1: Load and validate environment variables**
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	keycloakRealm := os.Getenv("KEYCLOAK_REALM")
	keycloakClientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	// Validate that all required variables are set
	if keycloakURL == "" || keycloakRealm == "" || keycloakClientID == "" || keycloakSecret == "" {
		return 0, fmt.Errorf("missing one or more required environment variables: KEYCLOAK_URL, KEYCLOAK_REALM, KEYCLOAK_CLIENT_ID, KEYCLOAK_CLIENT_SECRET")
	}

	// **Step 2: Construct the logout URL**
	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", keycloakURL, keycloakRealm)

	// **Step 3: Build form data using variables**
	data := url.Values{}
	data.Set("client_id", keycloakClientID)
	data.Set("client_secret", keycloakSecret)
	data.Set("refresh_token", refreshToken)

	// **Step 4: Create the HTTP request**
	req, err := http.NewRequestWithContext(ctx, "POST", logoutURL, strings.NewReader(data.Encode()))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// **Step 5: Set up the HTTP client**
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// **Step 6: Execute the request**
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute logout request: %w", err)
	}
	defer resp.Body.Close()

	// **Step 7: Handle non-200 responses**
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to logout from Keycloak: %s, response: %s", resp.Status, string(body))
	}

	// **Step 8: Return the status code on successful logout**
	return true, nil
}
