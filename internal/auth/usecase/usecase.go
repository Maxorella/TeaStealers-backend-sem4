package usecase

import (
	"context"
	"errors"
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
		Name:         data.Name,
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

func (u *AuthUsecase) GetUserByID(ctx context.Context, uID uuid.UUID) (*models.User, error) {
	user, err := u.repo.GetUserByID(ctx, uID)
	return user, err
}

func (u *AuthUsecase) UpdateUserPassword(data *models.UserUpdatePassword) (string, time.Time, error) {
	oldPasswordHash := utils.GenerateHashString(data.OldPassword)
	newPasswordHash := utils.GenerateHashString(data.NewPassword)
	if oldPasswordHash == newPasswordHash {
		return "", time.Now(), errors.New("passwords must not match")
	}
	if err := u.repo.CheckUserPassword(data.ID, oldPasswordHash); err != nil {
		return "", time.Now(), errors.New("invalid old password")
	}
	level, err := u.repo.UpdateUserPassword(data.ID, newPasswordHash)
	if err != nil {
		return "", time.Now(), errors.New("incorrect id or passwordhash")
	}
	user := &models.User{
		ID:          data.ID,
		LevelUpdate: level,
	}
	token, exp, err := jwt.GenerateToken(user)
	if err != nil {
		return "", time.Now(), errors.New("unable to generate token")
	}
	return token, exp, nil
}
