// routes package
package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Router(app *echo.Echo) {
	// Status route
	app.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}
