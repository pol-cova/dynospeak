// routes package
package routes

import (
	"backend/internal/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Router(app *echo.Echo) {
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
}
