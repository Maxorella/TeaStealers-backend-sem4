//go:generate mockgen -destination=mock/${GOFILE} -package=${GOPACKAGE}_mock -source=${GOFILE}
package auth

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"time"

	"github.com/satori/uuid"
)

// AuthUsecase represents the usecase interface for authentication.
type AuthUsecase interface {
	SignUp(context.Context, *models.UserSignUpData) (*models.User, string, time.Time, error)
	Login(context.Context, *models.UserLoginData) (*models.User, string, time.Time, error)
	CheckAuth(context.Context, string) (uuid.UUID, error)
	UpdateUserPassword(*models.UserUpdatePassword) (string, time.Time, error) // тут менять левел юзера + генерировать новый жвт
	// CheckUserPassword(uuid.UUID, string) error
	GetUserByID(context.Context, uuid.UUID) (*models.User, error)
}

// AuthRepo represents the repository interface for authentication.
type AuthRepo interface {
	CreateUser(ctx context.Context, newUser *models.User) (*models.User, error)
	CheckUser(ctx context.Context, login string, passwordHash string) (*models.User, error)
	GetUserByID(ctx context.Context, uID uuid.UUID) (*models.User, error)
	GetUserLevelById(id uuid.UUID) (int, error)
	UpdateUserPassword(uuid.UUID, string) (int, error)
	CheckUserPassword(uuid.UUID, string) error
}
