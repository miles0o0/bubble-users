package database

import (
	"context"

	"github.com/miles0o0/bubble-users/graph/model"
)

// Repository defines the database operations for the service.
type Repository interface {
	// User Operations
	GetUserData(ctx context.Context, userID string) (*model.User, error)
	GetFriends(ctx context.Context, userID string) ([]*model.User, error)
	GetSettings(ctx context.Context, userID string) (*model.Settings, error)
	SetSettings(ctx context.Context, userID string, settings *model.SettingsInput) (*model.Settings, error)

	// Message Operations
	GetDMs(ctx context.Context, userID string, friendID *string) ([]*model.Message, error)
}
