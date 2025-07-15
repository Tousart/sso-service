package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tousart/sso/config"
)

type TokenRepo struct {
	db *redis.Client
}

func NewTokenRepo(cfg *config.Config) (*TokenRepo, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB_ID,
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
