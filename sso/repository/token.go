package repository

import (
	"context"
	"time"
)

type TokenRepo interface {
	SetToken(ctx context.Context, key string, value string, duration time.Duration) error
}
