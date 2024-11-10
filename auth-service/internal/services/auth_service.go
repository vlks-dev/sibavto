package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vlks-dev/sibavto/auth-service/internal/domain"
	"github.com/vlks-dev/sibavto/auth-service/internal/storage"
	"github.com/vlks-dev/sibavto/shared/utils/hash"
	"github.com/vlks-dev/sibavto/shared/utils/token"
	"log/slog"
	"time"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type AuthService struct {
	authRepository storage.AuthRepository
	log            *slog.Logger
	tokenService   *token.JWTService
}

func NewAuthService(repository storage.AuthRepository, logger *slog.Logger, tokenService *token.JWTService) *AuthService {
	return &AuthService{
		authRepository: repository,
		log:            logger,
		tokenService:   tokenService,
	}
}

func (s *AuthService) Register(ctx context.Context, req *domain.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if err := utils.Validate.Struct(req); err != nil {
		return uuid.Nil, err
	}
	_, err := s.authRepository.GetUserByEmail(ctx, req.Email)
	if err == nil {
		s.log.Info("User already exists", "email", req.Email)
		return uuid.Nil, ErrUserAlreadyExists
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.log.Error(err.Error())
		return uuid.Nil, err
	}

	var reqCopy = *req
	s.log.Debug("creating new user", "request", reqCopy)
	hashCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			errCh <- err
			return
		}
		hashCh <- hashedPassword
	}()

	select {
	case hashedPassword := <-hashCh:
		reqCopy.Password = hashedPassword
	case err := <-errCh:
		return uuid.Nil, err
	case <-ctx.Done():
		return uuid.Nil, ctx.Err()
	}

	return s.authRepository.CreateUser(ctx, &reqCopy)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		s.log.Info("User login failed", "email", email)
		s.log.Debug(
			"User login failed",
			"email",
			email,
			"password",
			password,
			"db_password",
			user.Password,
		)
		return "", err
	}

	token, err := s.tokenService.CreateJWT(user.ID, user.Roles)
	if err != nil {
		s.log.Error("Failed to create JWT", "error", err.Error())
		return "", errors.New("failed to create JWT")
	}

	/*	redisKey := fmt.Sprintf("user_token:%s", user.ID)
		err = s.redisClient.Set(
			ctx,
			redisKey,
			token,
			time.Hour*time.Duration(12),
		).Err()
		if err != nil {
			s.log.Error("failed to save token to redis", err.Error())
			return "", err
		}*/

	return token, nil
}
