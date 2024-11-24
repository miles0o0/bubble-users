package database

import (
	"context"
	"database/sql"

	"github.com/miles0o0/bubble-users/graph/model"
)

type PostgresRepository struct {
	DB *sql.DB
}

// User Operations
func (r *PostgresRepository) GetUserData(ctx context.Context, userID string) (*model.User, error) {
	query := `SELECT id, name, username, email FROM users WHERE id = $1`
	user := &model.User{}
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) GetFriends(ctx context.Context, userID string) ([]*model.User, error) {
	query := `SELECT u.id, u.name, u.username, u.email 
              FROM users u
              JOIN friends f ON f.friend_id = u.id
              WHERE f.user_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		friends = append(friends, user)
	}
	return friends, nil
}

func (r *PostgresRepository) GetSettings(ctx context.Context, userID string) (*model.Settings, error) {
	query := `SELECT theme, notifications FROM settings WHERE user_id = $1`
	settings := &model.Settings{}
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&settings.Theme, &settings.Notifications)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *PostgresRepository) SetSettings(ctx context.Context, userID string, settings *model.SettingsInput) (*model.Settings, error) {
	query := `INSERT INTO settings (user_id, theme, notifications)
              VALUES ($1, $2, $3)
              ON CONFLICT (user_id)
              DO UPDATE SET theme = $2, notifications = $3`
	_, err := r.DB.ExecContext(ctx, query, userID, settings.Theme, settings.Notifications)
	if err != nil {
		return nil, err
	}
	return &model.Settings{
		Theme:         settings.Theme,
		Notifications: settings.Notifications,
	}, nil
}

// Message Operations
func (r *PostgresRepository) GetDMs(ctx context.Context, userID string, friendID *string) ([]*model.Message, error) {
	query := `SELECT id, sender_id, receiver_id, content, timestamp 
              FROM messages 
              WHERE (sender_id = $1 AND receiver_id = $2) 
                 OR (sender_id = $2 AND receiver_id = $1)`
	rows, err := r.DB.QueryContext(ctx, query, userID, friendID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
