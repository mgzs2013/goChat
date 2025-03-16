package repository

import (
	"database/sql"
	"fmt"
	"goChat/internal/database"

	_ "github.com/lib/pq"
)

// GetUserRole fetches the role name for a given user ID
func GetUserRole(userID int) (string, error) {
	var roleName string
	err := database.Pool.QueryRow(`
        SELECT r.name
        FROM users u
        JOIN roles r ON u.role_id = r.id
        WHERE u.id = $1
    `, userID).Scan(&roleName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user role: %w", err)
	}
	return roleName, nil
}

func CreateUser(Pool *sql.DB, username, roleName string) (int, error) {
	// Find role_id based on roleName
	var roleID int
	err := Pool.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch role ID: %w", err)
	}

	// Insert the new user
	var userID int
	err = Pool.QueryRow(
		"INSERT INTO users (username, role_id) VALUES ($1, $2) RETURNING id",
		username, roleID,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}
