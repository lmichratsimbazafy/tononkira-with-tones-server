package auth

import (
	"errors"
	"time"

	"lmich.com/tononkira/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	SecretKey     string
	TokenDuration time.Duration
}

// NewJWTService creates a new JWTService
func NewJWTService(secretKey string, tokenDuration time.Duration) *JWTService {
	return &JWTService{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration,
	}
}

// GenerateToken generates a JWT for a given user
func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.TokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.SecretKey))
}

// ValidateToken validates a JWT and returns the associated user
func (s *JWTService) ValidateToken(encodedToken string) (*domain.User, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID := claims["user_id"].(string)
	return &domain.User{ID: userID}, nil
}
