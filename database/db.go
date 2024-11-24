package database

import (
	"database/sql"
	"log"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt string
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	rows, err := r.DB.Query("SELECT id, name, email, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			log.Println("Error scanning user:", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}
