package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Pool *sql.DB

func InitDB(connectionString string) (*sql.DB, error) {
	var err error

	// Connect to the database
	Pool, err = sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = Pool.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established!")
	return Pool, nil
}
