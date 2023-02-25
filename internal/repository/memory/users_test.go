package memory

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"reflect"
	"sync"
	"testing"
)

func TestNewUsersRepository(t *testing.T) {
	repo := NewUsersRepository()
	require.NotNil(t, repo)

	assert.Len(t, repo.(*UsersRepository).storage, 0)
}

func TestUsersRepository_FindByLogin(t *testing.T) {
	type fields struct {
		storage map[string]models.User
	}
	type args struct {
		login string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.User
	}{
		{
			name:   "storage is empty",
			fields: fields{storage: make(map[string]models.User)},
			args:   args{login: "qwerty"},
			want:   nil,
		},
		{
			name: "no such user with such login",
			fields: fields{
				storage: map[string]models.User{
					"test": {
						Login:    "test",
						Password: "secret",
					},
				},
			},
			args: args{login: "qwerty"},
			want: nil,
		},
		{
			name: "user with such login exists",
			fields: fields{
				storage: map[string]models.User{
					"qwerty": {
						Login:    "qwerty",
						Password: "super_secret",
					},
					"test": {
						Login:    "test",
						Password: "secret",
					},
				},
			},
			args: args{login: "qwerty"},
			want: &models.User{
				Login:    "qwerty",
				Password: "super_secret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UsersRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			got, err := repo.FindByLogin(context.Background(), tt.args.login)
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsersRepository_Create(t *testing.T) {
	rand.Seed(1)

	type fields struct {
		storage map[string]models.User
	}
	type args struct {
		login      string
		password   string
		aesSecret  string
		privateKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "user with such login already exists",
			fields: fields{
				storage: map[string]models.User{
					"qwerty": {
						Login:    "qwerty",
						Password: "secret",
					},
				},
			},
			args: args{
				login:      "qwerty",
				password:   "super_secret",
				aesSecret:  "aes hashed secret",
				privateKey: "hashed private key",
			},
			wantErr: true,
		},
		{
			name: "user created successfully",
			fields: fields{
				storage: map[string]models.User{
					"test": {
						Login:    "test",
						Password: "secret",
					},
				},
			},
			args: args{
				login:      "qwerty",
				password:   "super_secret",
				aesSecret:  "aes hashed secret",
				privateKey: "hashed private key",
			},
			want: &models.User{
				Login:      "qwerty",
				Password:   "super_secret",
				AesSecret:  "aes hashed secret",
				PrivateKey: "hashed private key",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UsersRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}
			ctx := context.Background()

			if err := repo.Create(ctx, tt.args.login, tt.args.password, tt.args.aesSecret, tt.args.privateKey); (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			createdUser, err := repo.FindByLogin(ctx, tt.args.login)
			require.NoError(t, err)

			tt.want.ID = createdUser.ID
			if !reflect.DeepEqual(createdUser, tt.want) {
				t.Errorf("After creating a user, found by the same login and got = %v, want %v", createdUser, tt.want)
			}
		})
	}
}
