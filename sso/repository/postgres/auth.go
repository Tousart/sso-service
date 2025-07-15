package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tousart/sso/config"
	"github.com/tousart/sso/pkg"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *sql.DB
}

func CreateAuthRepo(cfg *config.Config) (*AuthRepo, error) {
	db, err := pkg.ConnectToDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %v", err)
	}

	return &AuthRepo{db: db}, nil
}

func (r *AuthRepo) Login(ctx context.Context, login string, password string) (string, error) {
	var (
		exists       bool
		hashPassword string
		userID       string
	)

	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", login).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to scan login: %v", err)
	}

	if !exists {
		return "", errors.New("user not exists")
	}

	err = r.db.QueryRow("SELECT hash_password, user_id FROM users WHERE login = $1", login).Scan(&hashPassword, &userID)
	if err != nil {
		return "", errors.New("failed to select password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return "", errors.New("incorrect password")
	}

	return userID, nil
}

func (r *AuthRepo) Register(ctx context.Context, login string, hashPassword []byte, email string, userID string) error {
	var exists bool

	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", login).Scan(&exists)
	if err != nil {
		return errors.New("failed to scan login")
	}

	if exists {
		return errors.New("user exists")
	}

	_, err = r.db.Exec("INSERT INTO users (user_id, login, hash_password, email) VALUES ($1, $2, $3, $4)", userID, login, string(hashPassword), email)
	if err != nil {
		return errors.New("user insertion error")
	}

	return nil
}
