package service

import (
	"context"
	"fmt"

	"github.com/PashakArt/file-server/internal/domain"
	"github.com/PashakArt/file-server/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(ur *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  ur,
		jwtSecret: []byte(jwtSecret),
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

	claims := jwt.MapClaims{
		"sub": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("AuthService:Login:token.SignedString - %w", err)
	}

	return signedToken, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	return nil
}
