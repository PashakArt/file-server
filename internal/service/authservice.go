package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PashakArt/file-server/internal/db/repository"
	"github.com/PashakArt/file-server/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	redisClient *redis.Client
	tokenTTL    time.Duration
}

func NewAuthService(ur *repository.UserRepository, redisClient *redis.Client, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:    ur,
		redisClient: redisClient,
		tokenTTL:    tokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, login, password string) error {
	existingUser, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return fmt.Errorf("AuthService:Register:userRepo.GetByLogin - %w", err)
	}

	if existingUser != nil {
		return domain.ErrUserAlreadyExists
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("AuthService:Register:bcrypt.GenerateFromPassword - %w", err)
	}

	_, err = s.userRepo.Create(ctx, login, string(hashPassword))
	if err != nil {
		return fmt.Errorf("AuthService:Register:userRepo.Create - %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("AuthService:Login:userRepo.GetByLogin - %w", err)
	}

	if user == nil {
		return "", domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token := uuid.New().String()
	redisKey := fmt.Sprintf("token:%s", token)
	err = s.redisClient.Set(ctx, redisKey, user.ID, s.tokenTTL).Err()
	if err != nil {
		return "", fmt.Errorf("AuthService:Login:redisClient.Set - %w", err)
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	return nil
}
