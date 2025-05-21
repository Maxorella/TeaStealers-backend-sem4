package models

import (
	"github.com/satori/uuid"
)

// UserSignUpData represents user information for signup
type UserSignUpData struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// UserLoginData represents user information for login
type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	AuthToken string `json:"token"`
}

// User represents user information
type User struct {
	ID           uuid.UUID `json:"id,omitempty"`
	PasswordHash string    `json:"-"`
	LevelUpdate  int       `json:"-"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name,omitempty"`
	IsDeleted    bool      `json:"-"`
	Token        string    `json:"token,omitempty"`
}

type UserUpdatePassword struct {
	// ID uniquely identifies the user.
	ID uuid.UUID `json:"id"`
	// OldPassword ...
	OldPassword string `json:"oldPassword"`
	// NewPassword ...
	NewPassword string `json:"newPassword"`
}
