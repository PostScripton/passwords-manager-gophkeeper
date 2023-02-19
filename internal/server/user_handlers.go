package server

import (
	"context"
	"errors"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/services"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errInternalServerError = status.Error(codes.Unknown, "Internal server error")

type UserServer struct {
	services *services.Services
	pb.UnimplementedUserServer
}

var _ pb.UserServer = (*UserServer)(nil)

func (u *UserServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, err := u.services.User.Create(ctx, in.Login, in.Password, in.AesSecret, in.PrivateKey)
	if err != nil {
		if errors.Is(err, repository.ErrLoginTaken) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		log.Error().Err(err).Msg("Creating user on register")
		return nil, errInternalServerError
	}

	token, err := u.services.Auth.LoginByUser(user)
	if err != nil {
		log.Error().Err(err).Int("user-id", user.ID).Msg("Authorizing user")
		return nil, errInternalServerError
	}

	log.Debug().Str("login", in.Login).Msg("User created successfully")

	return &pb.AuthResponse{
		Token:      token,
		AesSecret:  user.AesSecret,
		PrivateKey: user.PrivateKey,
	}, nil
}

func (u *UserServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := u.services.Auth.Login(ctx, in.Login, in.Password)
	if err != nil {
		if errors.Is(err, services.ErrCredentials) {
			return nil, status.Error(codes.Unauthenticated, "The credentials don't match any of our records")
		}

		log.Error().Err(err).Msg("Authorizing by credentials")
		return nil, errInternalServerError
	}

	user, err := u.services.User.FindByLogin(ctx, in.Login)
	if err != nil {
		return nil, errInternalServerError
	}

	return &pb.AuthResponse{
		Token:      token,
		AesSecret:  user.AesSecret,
		PrivateKey: user.PrivateKey,
	}, nil
}
