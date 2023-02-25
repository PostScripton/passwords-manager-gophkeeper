package memory

import (
	"context"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
)

func TestCredsSecretsRepository_Create(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	type args struct {
		userID         int
		website        string
		login          string
		encPassword    string
		additionalData string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "the creds already exist",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1234: {
						ID:             rand.Int63(),
						UID:            1234,
						Website:        "https://example.com",
						Login:          "qwerty",
						Password:       "encrypted_password",
						AdditionalData: "{}",
						UserID:         1,
					},
				},
			},
			args: args{
				userID:         1,
				website:        "https://example.com",
				login:          "qwerty",
				encPassword:    "encrypted_password",
				additionalData: "{}",
			},
			wantErr: assert.Error,
		},
		{
			name:   "creds stored successfully",
			fields: fields{storage: map[int64]models.CredsSecret{}},
			args: args{
				userID:         1,
				website:        "https://example.com",
				login:          "qwerty",
				encPassword:    "encrypted_password",
				additionalData: "{}",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			err := repo.Create(context.Background(), tt.args.userID, tt.args.website, tt.args.login, tt.args.encPassword, tt.args.additionalData)
			if !tt.wantErr(t, err) {
				return
			}

			var found *models.CredsSecret
			for _, creds := range repo.storage {
				if creds.UserID == tt.args.userID &&
					creds.Website == tt.args.website &&
					creds.Login == tt.args.login &&
					creds.Password == tt.args.encPassword &&
					creds.AdditionalData == tt.args.additionalData {
					found = &creds
					break
				}
			}

			if found == nil {
				t.Error("created creds secret not found in the storage")
			}
		})
	}
}

func TestCredsSecretsRepository_GetById(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	type args struct {
		uid int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.CredsSecret
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty storage",
			fields:  fields{storage: map[int64]models.CredsSecret{}},
			args:    args{uid: 1234},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "found creds secret",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1234: {
						ID:       4321,
						UID:      1234,
						Website:  "https://example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
				},
			},
			args: args{uid: 1234},
			want: &models.CredsSecret{
				ID:       4321,
				UID:      1234,
				Website:  "https://example.com",
				Login:    "qwerty",
				Password: "encrypted_password",
				UserID:   1,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			got, err := repo.GetById(context.Background(), tt.args.uid)
			if !tt.wantErr(t, err, fmt.Sprintf("GetById(%v)", tt.args.uid)) {
				return
			}

			assert.Equalf(t, tt.want, got, "GetById(%v)", tt.args.uid)
		})
	}
}

func TestCredsSecretsRepository_Delete(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	type args struct {
		uid int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "empty storage",
			fields: fields{storage: map[int64]models.CredsSecret{}},
			args:   args{uid: 1234},
		},
		{
			name: "delete creds secret",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1234: {
						ID:       4321,
						UID:      1234,
						Website:  "https://example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
				},
			},
			args: args{uid: 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			err := repo.Delete(context.Background(), tt.args.uid)
			require.NoError(t, err)

			_, ok := repo.storage[tt.args.uid]
			assert.False(t, ok)
		})
	}
}

func TestCredsSecretsRepository_GetList(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	type args struct {
		userID int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int64
	}{
		{
			name:   "empty storage",
			fields: fields{storage: map[int64]models.CredsSecret{}},
			args:   args{userID: 1},
			want:   []int64{},
		},
		{
			name: "get user's elements out of many elements of other users",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1: {
						ID:       1,
						UID:      1,
						Website:  "https://b_example.com",
						Login:    "a_qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
					2: {
						ID:       2,
						UID:      2,
						Website:  "https://a_example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   3,
					},
					3: {
						ID:       3,
						UID:      3,
						Website:  "https://a_example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
					4: {
						ID:       4,
						UID:      4,
						Website:  "https://b_example.com",
						Login:    "b_qwerty",
						Password: "encrypted_password",
						UserID:   2,
					},
					5: {
						ID:       5,
						UID:      5,
						Website:  "https://b_example.com",
						Login:    "b_qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
				},
			},
			args: args{userID: 1},
			want: []int64{3, 1, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			got, err := repo.GetList(context.Background(), tt.args.userID)
			require.NoError(t, err)

			require.Len(t, got, len(tt.want), "the length of elements in list must be the same as wanted")

			for i, uid := range tt.want {
				assert.Equalf(t, uid, got[i].UID, "the order of elements should be desc, website first and login second")
			}
		})
	}
}

func TestCredsSecretsRepository_SetList(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	type args struct {
		list []models.CredsSecret
	}
	type assertions struct {
		len int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		assertions assertions
	}{
		{
			name:   "set list on empty storage",
			fields: fields{storage: map[int64]models.CredsSecret{}},
			args: args{
				list: []models.CredsSecret{
					{
						ID:       1,
						UID:      1,
						Website:  "https://example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
					{
						ID:       2,
						UID:      2,
						Website:  "https://example.com",
						Login:    "qwerty2",
						Password: "encrypted_password",
						UserID:   2,
					},
				},
			},
			assertions: assertions{
				len: 2,
			},
		},
		{
			name: "set list of new elements on not empty storage",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1: {
						ID:       1,
						UID:      1,
						Website:  "https://example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
					2: {
						ID:       2,
						UID:      2,
						Website:  "https://example.com",
						Login:    "qwerty2",
						Password: "encrypted_password",
						UserID:   2,
					},
				},
			},
			args: args{
				list: []models.CredsSecret{
					{
						ID:       3,
						UID:      3,
						Website:  "https://example.com",
						Login:    "qwerty3",
						Password: "encrypted_password",
						UserID:   3,
					},
					{
						ID:       4,
						UID:      4,
						Website:  "https://example.com",
						Login:    "qwerty4",
						Password: "encrypted_password",
						UserID:   4,
					},
				},
			},
			assertions: assertions{
				len: 4,
			},
		},
		{
			name: "set list and override some of elements",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1: {
						ID:       1,
						UID:      1,
						Website:  "https://example.com",
						Login:    "qwerty",
						Password: "encrypted_password",
						UserID:   1,
					},
					2: {
						ID:       2,
						UID:      2,
						Website:  "https://example.com",
						Login:    "qwerty2",
						Password: "encrypted_password",
						UserID:   2,
					},
				},
			},
			args: args{
				list: []models.CredsSecret{
					{
						ID:       2,
						UID:      2,
						Website:  "https://example.com",
						Login:    "new login",
						Password: "encrypted_password2",
						UserID:   2,
					},
					{
						ID:       3,
						UID:      3,
						Website:  "https://example.com",
						Login:    "qwerty3",
						Password: "encrypted_password",
						UserID:   3,
					},
				},
			},
			assertions: assertions{
				len: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			err := repo.SetList(context.Background(), tt.args.list)
			require.NoError(t, err)

			assert.Len(t, repo.storage, tt.assertions.len)
		})
	}
}

func TestCredsSecretsRepository_Truncate(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "empty storage",
			fields: fields{storage: map[int64]models.CredsSecret{}},
		},
		{
			name: "delete all the elements",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1: {
						ID:      rand.Int63(),
						UID:     1,
						Website: "https://example.com",
						Login:   "qwerty",
						UserID:  1,
					},
					2: {
						ID:      rand.Int63(),
						UID:     2,
						Website: "https://example.com",
						Login:   "qwerty1234",
						UserID:  1,
					},
					3: {
						ID:      rand.Int63(),
						UID:     3,
						Website: "https://example2.com",
						Login:   "zxc",
						UserID:  2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			err := repo.Truncate(context.Background())
			require.NoError(t, err)

			assert.Len(t, repo.storage, 0, "after truncating, the storage should be empty")
		})
	}
}

func TestCredsSecretsRepository_checkCredsSecretExists(t *testing.T) {
	type fields struct {
		storage map[int64]models.CredsSecret
	}
	type args struct {
		userID  int
		website string
		login   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "empty storage",
			fields: fields{storage: map[int64]models.CredsSecret{}},
			args: args{
				userID:  1,
				website: "https://example.com",
				login:   "qwerty",
			},
			want: false,
		},
		{
			name: "no exact login for this website",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1: {
						ID:      1,
						UID:     rand.Int63(),
						Website: "https://example.com",
						Login:   "some_login",
						UserID:  1,
					},
					2: {
						ID:      2,
						UID:     rand.Int63(),
						Website: "https://example2.com",
						Login:   "qwerty",
						UserID:  1,
					},
				},
			},
			args: args{
				userID:  1,
				website: "https://example.com",
				login:   "qwerty",
			},
			want: false,
		},
		{
			name: "login exists",
			fields: fields{
				storage: map[int64]models.CredsSecret{
					1: {
						ID:      1,
						UID:     rand.Int63(),
						Website: "https://example.com",
						Login:   "some_login",
						UserID:  1,
					},
					2: {
						ID:      2,
						UID:     rand.Int63(),
						Website: "https://example.com",
						Login:   "qwerty",
						UserID:  1,
					},
				},
			},
			args: args{
				userID:  1,
				website: "https://example.com",
				login:   "qwerty",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CredsSecretsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			got := repo.checkCredsSecretExists(tt.args.userID, tt.args.website, tt.args.login)

			assert.Equalf(t, tt.want, got, "checkCredsSecretExists(%v, %v, %v)", tt.args.userID, tt.args.website, tt.args.login)
		})
	}
}
