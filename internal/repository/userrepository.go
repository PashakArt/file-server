package repository

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/PashakArt/file-server/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed query/get_user_by_login_query.sql
	getByLoginQuery string

	//go:embed query/create_user_query.sql
	createQuery string
)

type UserRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		dbPool: dbPool,
	}
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	err := r.dbPool.QueryRow(ctx, getByLoginQuery, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository:GetByLogin - %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, login, passwordHash string) (*uuid.UUID, error) {
	userID := uuid.New()

	_, err := r.dbPool.Exec(ctx, createQuery, userID, login, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("UserRepository:Create - %w", err)
	}

	return &userID, nil
}
