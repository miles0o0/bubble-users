package util

import (
	"context"

	"github.com/miles0o0/bubble-users/graph/model"
)

// UserLogin handles the login logic for Keycloak
func UserLogin(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	// Delegate login to the Keycloak-specific function
	return keycloakLogin(ctx, username, password)
}
