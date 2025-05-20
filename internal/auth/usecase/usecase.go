package usecase

import (
	"context"
	"github.com/TeaStealers-backend-sem4/internal/auth"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/jwt"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"time"

	"github.com/satori/uuid"
)

// AuthUsecase represents the usecase for authentication.
type AuthUsecase struct {
	repo auth.AuthRepo
}

// NewAuthUsecase creates a new instance of AuthUsecase.
func NewAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

// SignUp handles the user registration process.
func (u *AuthUsecase) SignUp(ctx context.Context, data *models.UserSignUpData) (*models.User, string, time.Time, error) {
	newUser := &models.User{
		ID:           uuid.NewV4(),
		Email:        data.Email,
		PasswordHash: utils.GenerateHashString(data.Password),
	}

	userResponse, err := u.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, "", time.Now(), err
	}

	token, exp, err := jwt.GenerateToken(newUser)
	if err != nil {
		return nil, "", time.Now(), err
	}

	return userResponse, token, exp, nil
}

// Login handles the user login process.
func (u *AuthUsecase) Login(ctx context.Context, data *models.UserLoginData) (*models.User, string, time.Time, error) {
	user, err := u.repo.CheckUser(ctx, data.Email, utils.GenerateHashString(data.Password))
	if err != nil {
		return nil, "", time.Now(), err
	}

	token, exp, err := jwt.GenerateToken(user)
	if err != nil {
		return nil, "", time.Now(), err
	}

	return user, token, exp, nil
}

// CheckAuth checking autorizing
func (u *AuthUsecase) CheckAuth(ctx context.Context, token string) (uuid.UUID, error) {
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return uuid.Nil, err
	}
	id, _, err := jwt.ParseClaims(claims)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
