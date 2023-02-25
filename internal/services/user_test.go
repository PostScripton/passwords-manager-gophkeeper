package services

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	mock_repository "github.com/PostScripton/passwords-manager-gophkeeper/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestUserService_FindByLogin(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockUsers
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		prepare func(f *fields, a args, user *models.User)
		args    args
		want    *models.User
	}{
		{
			name: "empty storage",
			prepare: func(f *fields, a args, user *models.User) {
				f.repo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(user, nil)
			},
			args: args{login: "qwerty"},
			want: nil,
		},
		{
			name: "found user",
			prepare: func(f *fields, a args, user *models.User) {
				f.repo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(user, nil)
			},
			args: args{login: "qwerty"},
			want: &models.User{
				ID:         1,
				Login:      "qwerty",
				Password:   "encrypted_password",
				AesSecret:  "encrypted_aes_secret",
				PrivateKey: "encrypted_private_key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockUsers(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f, tt.args, tt.want)
			}
			service := NewUserService(f.repo)

			got, err := service.FindByLogin(context.Background(), tt.args.login)
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Create(t *testing.T) {
	type fields struct {
		repo *mock_repository.MockUsers
	}
	type args struct {
		login          string
		hashedPassword string
		aesSecret      string
		privateKey     string
	}
	tests := []struct {
		name    string
		prepare func(f *fields, a args, user *models.User)
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "login is already taken",
			prepare: func(f *fields, a args, _ *models.User) {
				f.repo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(nil, repository.ErrLoginTaken)
			},
			args: args{
				login: "qwerty",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "created user successfully",
			prepare: func(f *fields, a args, user *models.User) {
				gomock.InOrder(
					f.repo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(nil, nil),
					f.repo.EXPECT().Create(gomock.Any(), a.login, a.hashedPassword, a.aesSecret, a.privateKey).Return(nil),
					f.repo.EXPECT().FindByLogin(gomock.Any(), a.login).Return(user, nil),
				)
			},
			args: args{
				login:          "qwerty",
				hashedPassword: "encrypted_password",
				aesSecret:      "encrypted_aes_secret",
				privateKey:     "encrypted_private_key",
			},
			want: &models.User{
				ID:         1,
				Login:      "qwerty",
				Password:   "encrypted_password",
				AesSecret:  "encrypted_aes_secret",
				PrivateKey: "encrypted_private_key",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repository.NewMockUsers(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f, tt.args, tt.want)
			}
			s := NewUserService(f.repo)

			got, err := s.Create(context.Background(), tt.args.login, tt.args.hashedPassword, tt.args.aesSecret, tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
