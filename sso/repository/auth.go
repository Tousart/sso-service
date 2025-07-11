package repository

import "context"

type AuthRepo interface {
	Login(ctx context.Context, login string, password string) (string, error)
	Register(ctx context.Context, login string, hashPassword []byte, email string, userID string) error
}
