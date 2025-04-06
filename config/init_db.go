package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "filecheck.db"
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Database file does not exist, create and initialize it
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}
		file.Close()
		// log.Println("Database file created successfully.")
	} else {
		// log.Println("Database file already exists.")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Initialize database schema if necessary
	if err := initializeSchema(db); err != nil {
		return nil, err
	}

	// log.Println("Database initialized successfully.")
	return db, nil
}

func initializeSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
    	password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL,
		hash VARCHAR(64) NOT NULL,
		uploaded_at TIMESTAMP NOT NULL,
		location VARCHAR(1024) NOT NULL
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	fmt.Println("Database schema initialized successfully!")
	return nil
}
