package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/miles0o0/bubble-users/graph/model"

	"github.com/Nerzal/gocloak/v13"
)

func keycloakLogin(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	// Load configuration
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	keycloakClientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloakSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	keycloakRealm := os.Getenv("KEYCLOAK_REALM")

	// Validate required variables
	if keycloakURL == "" || keycloakRealm == "" {
		return nil, fmt.Errorf("keycloak URL or realm is not configured")
	}

	// Initialize gocloak client
	client := gocloak.NewClient(keycloakURL)

	// Perform login
	token, err := client.Login(ctx, keycloakClientID, keycloakSecret, keycloakRealm, username, password)
	if err != nil {
		return nil, fmt.Errorf("keycloak login failed: %w", err)
	}

	// Marshal the JWT object into JSON
	tokenJSON, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JWT: %v", err)
	} else {
		log.Printf("Full JWT response: %s", string(tokenJSON))
	}

	// Map the response to your LoginResponse model
	return &model.LoginResponse{
		Token:   token.AccessToken,
		Refresh: token.RefreshToken,
	}, nil
}
