package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/satori/uuid"
)

type AuthRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	insert := `INSERT INTO users (id, email, name, passwordhash) VALUES ($1, $2, $3, $4)`
	if _, err := r.db.ExecContext(ctx, insert, user.ID, user.Email, user.Name, user.PasswordHash); err != nil {
		return nil, err
	}
	query := `SELECT id, email, passwordhash, levelupdate FROM users WHERE id = $1`

	res := r.db.QueryRow(query, user.ID)

	newUser := &models.User{}
	if err := res.Scan(&newUser.ID, &newUser.Email, &newUser.PasswordHash, &newUser.LevelUpdate); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (r *AuthRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, email, passwordhash, levelupdate FROM users WHERE email = $1`

	res := r.db.QueryRowContext(ctx, query, login)

	user := &models.User{}
	if err := res.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.LevelUpdate); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, name, passwordhash, levelupdate FROM users WHERE id = $1`

	res := r.db.QueryRowContext(ctx, query, id)

	user := &models.User{}
	if err := res.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.LevelUpdate); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepo) CheckUser(ctx context.Context, login string, passwordHash string) (*models.User, error) {
	user, err := r.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if user.PasswordHash != passwordHash {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func (r *AuthRepo) GetUserLevelById(id uuid.UUID) (int, error) {
	query := `SELECT levelupdate FROM users WHERE id = $1`

	res := r.db.QueryRow(query, id)

	level := 0
	if err := res.Scan(&level); err != nil {
		return 0, err
	}
	return level, nil
}

func (r *AuthRepo) UpdateUserPassword(id uuid.UUID, newPasswordHash string) (int, error) {
	query := `UPDATE users SET passwordhash=$1, levelupdate = levelupdate+1 WHERE id = $2`
	if _, err := r.db.Exec(query, newPasswordHash, id); err != nil {
		return 0, err
	}
	querySelect := `SELECT levelupdate FROM users WHERE id = $1`
	level := 0
	res := r.db.QueryRow(querySelect, id)
	if err := res.Scan(&level); err != nil {
		return 0, err
	}
	return level, nil
}

func (r *AuthRepo) CheckUserPassword(id uuid.UUID, passwordHash string) error {
	passwordHashCur := ""
	querySelect := `SELECT passwordhash FROM users WHERE id = $1`
	res := r.db.QueryRow(querySelect, id)
	if err := res.Scan(&passwordHashCur); err != nil {
		return err
	}
	if passwordHashCur != passwordHash {
		return errors.New("passwords don't match")
	}
	return nil
}
