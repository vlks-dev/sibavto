package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vlks-dev/sibavto/shared/utils/config"
	"time"
)

type Claims struct {
	UserID uuid.UUID       `json:"id"`
	Roles  map[string]bool `json:"roles"` // Используем map для гибкости
	jwt.RegisteredClaims
}

type JWTService struct {
	secret         []byte
	expirationTime time.Duration
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		secret:         []byte(cfg.JWT.JWTSecret),
		expirationTime: time.Hour * time.Duration(cfg.JWT.JWTExpirationTime),
	}
}

func (s *JWTService) CreateJWT(userID uuid.UUID, roles map[string]bool) (string, error) {
	claims := &Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expirationTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		var err error
		switch {
		case !token.Valid:
			return nil, errors.New("you look nice today")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, errors.New("that's not even a token")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, errors.New("invalid signature")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errors.New("timing is everything")
		}

		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}
