package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tousart/sso/domain/models"
	"github.com/tousart/sso/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo  repository.AuthRepo
	tokenRepo repository.TokenRepo
	sender    repository.Sender
}

func CreateAuthService(authRepo repository.AuthRepo, tokenRepo repository.TokenRepo, sender repository.Sender) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		tokenRepo: tokenRepo,
		sender:    sender,
	}
}

func (s *AuthService) Login(ctx context.Context, login string, password string) (token string, err error) {
	userID, err := s.authRepo.Login(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to get userID: %v", err)
	}

	duration := time.Hour * 12

	token, err = generateToken(userID, login, duration)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	if err = s.tokenRepo.SetToken(context.Background(), userID, token, duration); err != nil {
		return "", fmt.Errorf("redis error (token has been deleted): %v", err)
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, login string, password string, email string) (err error) {
	hashPassword, err := hash(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	userID := generateID()

	err = s.authRepo.Register(ctx, login, hashPassword, email, userID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	emailMessage := models.EmailMessage{
		Login: login,
		Email: email,
	}

	message, err := json.Marshal(emailMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = s.sender.SendMessage(context.Background(), []byte("hello"), message)
	if err != nil {
		return fmt.Errorf("failed to send to broker: %v", err)
	}

	return nil
}

func hash(password string) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash: %v", err)
	}
	return hashPassword, nil
}

func generateID() string {
	id := uuid.NewString()
	return id
}

func generateToken(userID string, login string, duration time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    userID,
			"login": login,
			"exp":   time.Now().Add(duration).Unix(),
		})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}
