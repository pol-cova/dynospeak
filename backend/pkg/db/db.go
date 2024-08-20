package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./pkg/db/backend.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Create SQL tables
	createAuthTables()
}

func createAuthTables() {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                email TEXT NOT NULL UNIQUE,
                username TEXT NOT NULL UNIQUE,
                password TEXT NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	createBlackListTokensTable := `
    CREATE TABLE IF NOT EXISTS blacklist (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                token TEXT NOT NULL UNIQUE,
                logout_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                );`
	_, err = DB.Exec(createBlackListTokensTable)
	if err != nil {
		log.Fatalf("Could not create blacklist table: %v", err)
	}
}
