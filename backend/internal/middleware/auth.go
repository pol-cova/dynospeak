package middleware

import (
	"backend/internal/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized", "message": "token required"})
		}
		// check if token is not in blacklist
		isBlacklisted := auth.IsTokenBlacklisted(token)
		if isBlacklisted {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized", "message": "token expired, please login again"})
		}

		userId, username, err := auth.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized", "message": "invalid token"})
		}
		c.Set("userId", userId)
		c.Set("username", username)

		return next(c)
	}
}
