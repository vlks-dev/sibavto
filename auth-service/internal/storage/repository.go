package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlks-dev/sibavto/auth-service/internal/domain"
	"log/slog"
	"time"
)

//go:generate mockery --name=UserRepository --output=./mocks --case=underscore
type AuthRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type UserStorage struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewUserStorage(pool *pgxpool.Pool, log *slog.Logger) *UserStorage {
	return &UserStorage{
		pool: pool,
		log:  log,
	}
}

func (s *UserStorage) CreateUser(ctx context.Context, req *domain.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id uuid.UUID
	err := s.withTransaction(ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO users (name, surname,patronymic, email, hash)
		  VALUES ($1, $2, $3, $4, $5) RETURNING id`
		return tx.QueryRow(ctx, q, req.Name, req.Surname, req.Patronymic, req.Email, req.Password).Scan(&id)
	})
	if err != nil {
		s.log.Error("failed to register user", "error", err)
		return uuid.Nil, err
	}
	return id, nil
}

func (s *UserStorage) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user domain.User
	err := s.withTransaction(ctx, func(tx pgx.Tx) error {
		q := `SELECT id, name, surname, email FROM users WHERE id = $1`
		return tx.QueryRow(ctx, q, id).Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Email,
		)
	})
	if err != nil {
		s.log.Error("failed to get user by id", "error", err)
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user domain.User
	err := s.withTransaction(ctx, func(tx pgx.Tx) error {
		q := `SELECT id, name, surname, email, hash FROM users WHERE email = $1`
		return tx.QueryRow(ctx, q, email).Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Email,
			&user.Password,
		)
	})

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		s.log.Debug("creating new user", "?unique", true)
		return nil, err
	} else if err != nil {
		s.log.Error("failed to get user by email", "error", err.Error())
	}
	return &user, nil
}

func (s *UserStorage) withTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	start := time.Now()
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		s.log.Error("Failed to start transaction")
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			s.log.Error("Failed to rollback transaction", "error", rbErr)
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error("Failed to commit transaction", "error", err)
		return err
	}
	s.log.Debug("successfully committed transaction",
		"query_time (ms)", time.Since(start).Milliseconds())

	return nil
}
