package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tousart/sso/domain/models"
	"github.com/tousart/sso/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   repository.AuthRepo
	sender repository.Sender
}

func CreateAuthService(repo repository.AuthRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, login string, password string) (token string, err error) {
	userID, err := s.repo.Login(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to get userID: %v", err)
	}

	token, err = generateToken(userID, login)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// ОТПРАВКА ТОКЕНА В REDIS

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, login string, password string, email string) (err error) {
	hashPassword, err := hash(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	userID := generateID()

	err = s.repo.Register(ctx, login, hashPassword, email, userID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	var emailMessage models.EmailMessage

	message, err := json.Marshal(emailMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = s.sender.SendMessage(ctx, []byte("hello"), message)
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

func generateToken(userID string, login string) (string, error) {
	secretKey := ""

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    userID,
			"login": login,
			"exp":   time.Now().Add(time.Hour * 12).Unix(),
		})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}
