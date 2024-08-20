package main

import (
	"backend/internal/routes"
	"backend/pkg/db"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	uri := os.Getenv("MONGODB_URI")
	if port == "" {
		port = "8080"
	}

	// Define app instance
	app := echo.New()

	// Init db
	db.InitDB()
	dbName := "dynoapp"
	mongoInstance, err := db.InitializeDB(uri, dbName)
	if err != nil {
		panic(err)
	}
	defer mongoInstance.CloseDB()

	// Router
	routes.Router(app)
	// Start the server
	app.Logger.Fatal(app.Start("0.0.0.0:" + port))
}
