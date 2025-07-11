package usecase

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, login string, password string) (string, error)
	Register(ctx context.Context, login string, password string, email string) error
}
