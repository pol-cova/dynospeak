package models

import (
	"backend/internal/auth"
	"backend/pkg/db"
	"errors"
	"log"
	"regexp"
	"strings"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Username string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := `INSERT INTO users(email, username, password) VALUES(?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	hashedPassword, err := auth.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Email, u.Username, hashedPassword)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err
}

func (u *User) Authenticate() error {
	// Authenticate user logic
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)
	var RetrievedPassword string
	err := row.Scan(&u.ID, &RetrievedPassword)
	if err != nil {
		return err
	}
	isAuthenticated := auth.ValidatePasswordHash(u.Password, RetrievedPassword)
	if !isAuthenticated {
		return errors.New("invalid credentials")
	}
	return nil
}

func (u User) Profile() (User, error) {
	query := `SELECT email, username FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, u.ID)
	err := row.Scan(&u.Email, &u.Username)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (u User) Delete() error {
	query := `DELETE FROM users WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u User) AuthValidator() error {
	email := strings.TrimSpace(u.Email)
	password := strings.TrimSpace(u.Password)

	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	matchEmail, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, email)
	if err != nil {
		log.Printf("Error matching email regex: %v", err)
		return errors.New("error validating email")
	}
	if !matchEmail {
		return errors.New("invalid email")
	}

	// Validate password length
	if len(password) < 8 {
		return errors.New("password must have at least 8 characters")
	}

	// Validate password for at least one uppercase letter
	matchUppercase, err := regexp.MatchString(`[A-Z]`, password)
	if err != nil || !matchUppercase {
		return errors.New("password must have at least one uppercase letter")
	}

	// Validate password for at least one lowercase letter
	matchLowercase, err := regexp.MatchString(`[a-z]`, password)
	if err != nil || !matchLowercase {
		return errors.New("password must have at least one lowercase letter")
	}

	// Validate password for at least one digit
	matchDigit, err := regexp.MatchString(`\d`, password)
	if err != nil || !matchDigit {
		return errors.New("password must have at least one number")
	}

	return nil
}
