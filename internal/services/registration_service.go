package services

import (
	"fmt"
	"goChat/internal/database"

	"golang.org/x/crypto/bcrypt"
)

// store registered user
//  id
//  username
//  password_hash
//  role

// User represents the user model

// RegisterUser hashes the password and stores the user in the database
func RegisterUser(username, password string) error {

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store the user in the database

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err = database.Pool.Exec(query, username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to store user in database: %v", err)
	}
	return nil

}
