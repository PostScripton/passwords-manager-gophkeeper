package memory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestNewSettingsRepository(t *testing.T) {
	repo := NewSettingsRepository()
	require.NotNil(t, repo)

	assert.Len(t, repo.(*SettingsRepository).storage, 0)
}

func TestSettingsRepository_Get(t *testing.T) {
	type fields struct {
		storage map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		existed bool
	}{
		{
			name:    "empty storage",
			fields:  fields{storage: make(map[string]string)},
			args:    args{key: "key1"},
			want:    "",
			existed: false,
		},
		{
			name:    "no such key",
			fields:  fields{storage: map[string]string{"key": "value"}},
			args:    args{key: "key1"},
			want:    "",
			existed: false,
		},
		{
			name:    "key existed",
			fields:  fields{storage: map[string]string{"key": "value", "key1": "value1"}},
			args:    args{key: "key1"},
			want:    "value1",
			existed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SettingsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			got, gotExisted, err := repo.Get(context.Background(), tt.args.key)
			require.NoError(t, err)

			assert.Equalf(t, tt.want, got, "Get(%v) returned wrong value", tt.args.key)
			assert.Equalf(t, tt.existed, gotExisted, "Get(%v) returned wrong existence", tt.args.key)
		})
	}
}

func TestSettingsRepository_Set(t *testing.T) {
	type fields struct {
		storage map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		existed bool
	}{
		{
			name:    "set a new value",
			fields:  fields{storage: map[string]string{"key": "value"}},
			args:    args{key: "key1", value: "value1"},
			existed: false,
		},
		{
			name:    "override a key",
			fields:  fields{storage: map[string]string{"key": "value"}},
			args:    args{key: "key", value: "another value"},
			existed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SettingsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}
			ctx := context.Background()

			existed, err := repo.Set(ctx, tt.args.key, tt.args.value)
			require.NoError(t, err)

			assert.Equalf(t, tt.existed, existed, "Set(%v, %v)", tt.args.key, tt.args.value)

			got, existed, err := repo.Get(ctx, tt.args.key)
			require.NoError(t, err)

			assert.Equalf(t, tt.args.value, got, "Get(%v) must return a just set value", tt.args.key)
		})
	}
}

func TestSettingsRepository_Delete(t *testing.T) {
	type fields struct {
		storage map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		existed bool
	}{
		{
			name:    "key does not exist",
			fields:  fields{storage: map[string]string{"key": "value"}},
			args:    args{key: "key1"},
			existed: false,
		},
		{
			name:    "key existed",
			fields:  fields{storage: map[string]string{"key": "value", "key1": "value1"}},
			args:    args{key: "key1"},
			existed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SettingsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}
			ctx := context.Background()

			existed, err := repo.Delete(ctx, tt.args.key)
			require.NoError(t, err)

			assert.Equalf(t, tt.existed, existed, "Delete(%v)", tt.args.key)

			got, existed, err := repo.Get(ctx, tt.args.key)
			require.NoError(t, err)

			assert.Emptyf(t, got, "Get(%v) should return empty string on just deleted key", tt.args.key)
			assert.Falsef(t, existed, "Get(%v) should not exist because the key has just been deleted", tt.args.key)
		})
	}
}

func TestSettingsRepository_Truncate(t *testing.T) {
	type fields struct {
		storage map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "empty storage",
			fields: fields{map[string]string{}},
		},
		{
			name: "delete all items from storage",
			fields: fields{
				storage: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &SettingsRepository{
				storage: tt.fields.storage,
				mu:      &sync.RWMutex{},
			}

			err := repo.Truncate(context.Background())
			require.NoError(t, err)

			assert.Lenf(t, repo.storage, 0, "After truncating storage, there should not be any data")
		})
	}
}
