// routes package
package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/pkg/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Router(app *echo.Echo, dbInstance *db.MongoDB) {
	// Status route
	app.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.POST("/signup", handlers.Signup)
	auth.POST("/login", handlers.Login)
	auth.GET("/logout", handlers.Logout)
	auth.DELETE("/delete", handlers.DeleteAccount)

	// Authenticated user routes
	user := app.Group("/user")
	user.Use(middleware.AuthMiddleware)

	// User profile
	user.GET("/me", handlers.Profile)

	// Route to create a new chat room
	user.POST("/chatrooms/new", func(c echo.Context) error {
		return handlers.CreateNewRoom(c, dbInstance.ChatRooms)
	})

	// WebSocket route for chat messages
	app.GET("/ws/chat", func(c echo.Context) error {
		return handlers.HandleWebSocket(c, dbInstance.Messages)
	})

	// Route to get all messages in a specific room
	app.GET("/rooms/:room_name/messages", func(c echo.Context) error {
		return handlers.GetMessagesInRoom(c, dbInstance.Messages)
	})

	// Route to get all chat rooms
	app.GET("/rooms", func(c echo.Context) error {
		return handlers.GetAllRooms(c, dbInstance.ChatRooms)
	})
}
