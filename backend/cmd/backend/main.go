package main

import (
	"backend/internal/routes"
	"backend/pkg/db"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define app instance
	app := echo.New()
	// Init db
	db.InitDB()

	// Router
	routes.Router(app)
	// Start the server
	err := app.Start(":" + port)
	if err != nil {
		panic("could not start the server")
	}
}
