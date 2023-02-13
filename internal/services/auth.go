package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrCredentials = errors.New("credentials don't match")

type authService struct {
	userRepo repository.Users
	secret   string
}

var _ Auth = (*authService)(nil)

func NewAuthService(repo repository.Users, secret string) Auth {
	return &authService{
		userRepo: repo,
		secret:   secret,
	}
}

func (s *authService) GetSecret() string {
	return s.secret
}

func (s *authService) LoginByUser(user *models.User) (string, error) {
	return s.generateJWT(user)
}

func (s *authService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrCredentials
	}

	if err = s.checkPassword(user.Password, password); err != nil {
		return "", ErrCredentials
	}

	return s.generateJWT(user)
}

func (s *authService) ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
}

func (s *authService) GetIDFromJWT(token *jwt.Token) (int, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["sub"].(float64)

		return int(id), nil
	} else {
		return 0, fmt.Errorf("invalid jwt token")
	}
}

func (s *authService) generateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * 60 * time.Minute).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = user.ID

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) checkPassword(hashedPassword, providedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword)); err != nil {
		return err
	}

	return nil
}
