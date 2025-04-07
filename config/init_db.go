package config

import (
	"database/sql"
	"fmt"
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
	schemaPath := "/home/tnxl/FileCheck/db/migration/schema.sql"

	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	fmt.Println("Database schema initialized successfully from file.")
	return nil
}
