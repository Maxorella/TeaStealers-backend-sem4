package models

import (
	"github.com/satori/uuid"
)

// UserSignUpData represents user information for signup
type UserSignUpData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginData represents user information for login
type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents user information
type User struct {
	ID           uuid.UUID `json:"id"`
	PasswordHash string    `json:"-"`
	LevelUpdate  int       `json:"-"`
	Email        string    `json:"email"`
	IsDeleted    bool      `json:"-"`
}
