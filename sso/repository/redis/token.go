package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepo struct {
	db *redis.Client
}

func NewTokenRepo() (*TokenRepo, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "password",
		DB:       1,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &TokenRepo{
		db: db,
	}, nil
}

func (tr *TokenRepo) SetToken(ctx context.Context, key string, value string, duration time.Duration) error {
	log.Println("sending token...")
	if err := tr.db.Set(ctx, key, value, duration).Err(); err != nil {
		return fmt.Errorf("failed to set token to redis: %v", err)
	}
	log.Println("tokent has been sent to redis")
	return nil
}
