package main

import (
	"backend/internal/routes"
	"backend/pkg/db"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {

	// Get the port and MongoDB URI from the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatalf("MONGODB_URI is required")
	}

	if port == "" {
		port = "8080"
	}

	// Define app instance
	app := echo.New()

	// Allow cors
	app.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
			AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.OPTIONS, echo.DELETE},
		}))
	// Init db
	db.InitDB()
	dbName := "dynoapp"
	mongoInstance, err := db.InitializeDB(uri, dbName)
	if err != nil {
		panic(err)
	}
	defer mongoInstance.CloseDB()

	// Router
	routes.Router(app, mongoInstance)
	// Start the server
	app.Logger.Fatal(app.Start("0.0.0.0:" + port))
}
