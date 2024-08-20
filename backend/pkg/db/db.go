package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var DB *sql.DB

type MongoDB struct {
	Client    *mongo.Client
	ChatRooms *mongo.Collection
	Messages  *mongo.Collection
}

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

func InitializeDB(uri string, dbName string) (*MongoDB, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")

	// Select the database and collections
	database := client.Database(dbName)
	chatRoomsCollection := database.Collection("chatrooms")
	messagesCollection := database.Collection("messages")

	// Return the MongoDB instance
	return &MongoDB{
		Client:    client,
		ChatRooms: chatRoomsCollection,
		Messages:  messagesCollection,
	}, nil
}

// CloseDB closes the MongoDB connection.
func (db *MongoDB) CloseDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.Client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect MongoDB client: %v", err)
	}

	log.Println("Disconnected from MongoDB!")
	return nil
}
