package services

import (
	"context"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	mock_repository "github.com/PostScripton/passwords-manager-gophkeeper/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func Test_authService_Login(t *testing.T) {
	type fields struct {
		userRepo *mock_repository.MockUsers
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *fields, a args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no user with such login",
			args: args{
				login:    "qwerty",
				password: "some_password",
			},
			prepare: func(f *fields, a args) {
				f.userRepo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(nil, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "password for the user does not match",
			args: args{
				login:    "qwerty",
				password: "some_password",
			},
			prepare: func(f *fields, a args) {
				user := &models.User{
					ID:       1234,
					Login:    "qwerty",
					Password: "enc_password",
				}

				f.userRepo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(user, nil)
			},
			wantErr: assert.Error,
		},
		{
			name: "user is authorized",
			args: args{
				login:    "qwerty",
				password: "secret",
			},
			prepare: func(f *fields, a args) {
				user := &models.User{
					ID:       1234,
					Login:    "qwerty",
					Password: "$2a$08$aMSa62GGHKmnT0QPndKzWOo0TV59E/DVjNS1as3l4EISCnGUxFjfq",
				}

				f.userRepo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(user, nil)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				userRepo: mock_repository.NewMockUsers(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f, tt.args)
			}

			service := NewAuthService(f.userRepo, "super_secret")

			tokenString, err := service.Login(context.Background(), tt.args.login, tt.args.password)
			if !tt.wantErr(t, err, fmt.Sprintf("Login(%v, %v)", tt.args.login, tt.args.password)) || err != nil {
				return
			}

			parts := strings.Split(tokenString, ".")
			assert.Lenf(t, parts, 3, "the must be 3 parts in jwt")
		})
	}
}

func Test_authService_GenerateJWT(t *testing.T) {
	service := &authService{secret: "super_secret1234"}
	user := models.User{
		ID:       1234,
		Login:    "qwerty",
		Password: "some_password",
	}

	got, err := service.GenerateJWT(&user)
	require.NoError(t, err)

	parts := strings.Split(got, ".")
	assert.Lenf(t, parts, 3, "the must be 3 parts in jwt")
}

func Test_authService_ParseJWT(t *testing.T) {
	service := &authService{secret: "super_secret1234"}
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Nzc0MzI2MjYsImlhdCI6MTY3NzM0NjIyNiwic3ViIjoxMjM0fQ.3_l9epc5CQq1U3R3jTtmqrBX105UnuqUwYuE7PHQHEs"

	got, err := service.ParseJWT(tokenString)
	require.NoError(t, err)

	assert.NotNil(t, got)
}

func Test_authService_GetIDFromJWT(t *testing.T) {
	service := &authService{secret: "super_secret"}
	wantedID := 1234
	user := models.User{
		ID:       wantedID,
		Login:    "qwerty",
		Password: "some_password",
	}

	tokenString, err := service.GenerateJWT(&user)
	require.NoError(t, err)

	token, err := service.ParseJWT(tokenString)
	require.NoError(t, err)

	got, err := service.GetIDFromJWT(token)
	require.NoError(t, err)

	assert.Equal(t, wantedID, got)
}

func Test_authService_HashPassword(t *testing.T) {
	service := &authService{}
	password := "secret"

	got, err := service.HashPassword(password)
	require.NoError(t, err)

	if got == password {
		t.Errorf("passwords must not be the same")
	}
	if len(got) <= len(password) {
		t.Errorf("hashed password should be larger")
	}
}

func Test_authService_checkPassword(t *testing.T) {
	type args struct {
		hashedPassword   string
		providedPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "passwords match",
			args: args{
				hashedPassword:   "$2a$08$aMSa62GGHKmnT0QPndKzWOo0TV59E/DVjNS1as3l4EISCnGUxFjfq",
				providedPassword: "secret",
			},
			wantErr: false,
		},
		{
			name: "empty hash",
			args: args{
				hashedPassword:   "",
				providedPassword: "secret",
			},
			wantErr: true,
		},
		{
			name: "wrong password for the hashe",
			args: args{
				hashedPassword:   "$2a$08$aMSa62GGHKmnT0QPndKzWOo0TV59E/DVjNS1as3l4EISCnGUxFjfq",
				providedPassword: "secRet",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &authService{}

			if err := service.checkPassword(tt.args.hashedPassword, tt.args.providedPassword); (err != nil) != tt.wantErr {
				t.Errorf("checkPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
