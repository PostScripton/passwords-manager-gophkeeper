package services

import (
	"context"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	secret       string
	settingsRepo repository.Settings
}

var _ Auth = (*AuthService)(nil)

func NewAuthService(secret string, settingsRepo repository.Settings) Auth {
	return &AuthService{
		secret:       secret,
		settingsRepo: settingsRepo,
	}
}

func (s *AuthService) CheckAuthorized(ctx context.Context) (bool, error) {
	tokenString, existed, err := s.settingsRepo.Get(ctx, "jwt")
	if err != nil {
		return false, fmt.Errorf("getting JWT from setting: %w", err)
	}
	if !existed {
		return false, nil
	}

	token, err := s.parseJWT(tokenString)
	if err != nil {
		return false, fmt.Errorf("parse jwt: %w", err)
	}

	if !token.Valid {
		return false, fmt.Errorf("invalid jwt stored in settings")
	}

	return true, nil
}

func (s *AuthService) GetID(ctx context.Context) (int, error) {
	tokenString, _, err := s.settingsRepo.Get(ctx, "jwt")
	if err != nil {
		return 0, fmt.Errorf("getting JWT from setting: %w", err)
	}

	token, err := s.parseJWT(tokenString)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["sub"].(float64)

		return int(id), nil
	} else {
		return 0, fmt.Errorf("invalid jwt token")
	}
}

func (s *AuthService) parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
}
