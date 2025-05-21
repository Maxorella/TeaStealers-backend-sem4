package models

import (
	"github.com/satori/uuid"
)

type UserSignUpData struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	AuthToken string `json:"token"`
}

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
	ID          uuid.UUID `json:"id"`
	OldPassword string    `json:"oldPassword"`
	NewPassword string    `json:"newPassword"`
}
