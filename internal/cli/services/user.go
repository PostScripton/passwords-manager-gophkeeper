package services

import (
	"context"
	"fmt"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"strings"
)

var ErrLoginIsTaken = fmt.Errorf("this login is already taken")
var ErrLoggedInAlready = fmt.Errorf("the user is already logged in")
var ErrCredentialsDontMatch = fmt.Errorf("the credentials don't match any of our records")

type UserService struct {
	client pb.UserClient
	repos  *repository.Repository
}

var _ User = (*UserService)(nil)

func NewUserService(client pb.UserClient, repos *repository.Repository) User {
	return &UserService{
		client: client,
		repos:  repos,
	}
}

func (s *UserService) Register(ctx context.Context, login, password, aesSecret, privateKey string) error {
	request := &pb.RegisterRequest{
		Login:      login,
		Password:   password,
		AesSecret:  aesSecret,
		PrivateKey: privateKey,
	}

	response, err := s.client.Register(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "taken") {
			return ErrLoginIsTaken
		}

		return err
	}

	if err = s.rememberMe(ctx, response.Token, response.AesSecret, response.PrivateKey); err != nil {
		return fmt.Errorf("register remember me: %w", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, login, password string) error {
	request := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	response, err := s.client.Login(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "credentials") {
			return ErrCredentialsDontMatch
		}

		return err
	}

	if err = s.rememberMe(ctx, response.Token, response.AesSecret, response.PrivateKey); err != nil {
		return fmt.Errorf("login remember me: %w", err)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context) error {
	if err := s.repos.Settings.Truncate(ctx); err != nil {
		return err
	}
	if err := s.repos.CredsSecrets.Truncate(ctx); err != nil {
		return err
	}

	return nil
}

func (s *UserService) rememberMe(ctx context.Context, jwt, aesSecret, privateKey string) error {
	if _, err := s.repos.Settings.Set(ctx, "jwt", jwt); err != nil {
		return err
	}

	if _, err := s.repos.Settings.Set(ctx, "aes_secret", aesSecret); err != nil {
		return err
	}

	if _, err := s.repos.Settings.Set(ctx, "private_key", privateKey); err != nil {
		return err
	}

	return nil
}
