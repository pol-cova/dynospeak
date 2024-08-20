package handlers

import (
	"backend/internal/auth"
	"backend/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Signup(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "invalid request"})
	}
	// Validate information
	err = user.AuthValidator()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "something went wrong"})
	}

	// Save user
	err = user.Save()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error(), "message": "could not create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "user created"})
}

func Login(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "message": "invalid request"})
	}

	err = user.Authenticate()
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error(), "message": "invalid credentials"})
	}
	token, err := auth.GenerateToken(user.Email, user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error(), "message": "could not generate token"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token, "message": "login successful"})
}

func Profile(c echo.Context) error {
	userId := c.Get("userId").(int64)
	user := models.User{ID: userId}
	user, err := user.Profile()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error(), "message": "could not get user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"email": user.Email, "username": user.Username})
}

func DeleteAccount(c echo.Context) error {
	userId := c.Get("userId").(int64)
	user := models.User{ID: userId}
	err := user.Delete()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error(), "message": "could not delete user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}

func Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	err := auth.LogoutToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error(), "message": "could not logout user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "user logged out"})
}
