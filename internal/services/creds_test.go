package services

import (
	"context"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	mock_repository "github.com/PostScripton/passwords-manager-gophkeeper/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCredsService_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockCredsSecrets(ctrl)
	userID := 1234
	expectedList := []*models.CredsSecret{
		{
			ID:       1,
			UID:      1,
			Website:  "https://example.com",
			Login:    "qwerty",
			Password: "enc_password",
			UserID:   1234,
		},
		{
			ID:       2,
			UID:      2,
			Website:  "https://example2.com",
			Login:    "zxc",
			Password: "enc_password2",
			UserID:   1234,
		},
	}
	repo.EXPECT().GetList(gomock.Any(), userID).Return(expectedList, nil)

	service := NewCredsService(repo)

	got, err := service.GetList(context.Background(), userID)
	require.NoError(t, err)

	assert.Equal(t, expectedList, got)
}

func TestCredsService_SetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock_repository.NewMockCredsSecrets(ctrl)
	list := []models.CredsSecret{
		{
			ID:       1,
			UID:      1,
			Website:  "https://example.com",
			Login:    "qwerty",
			Password: "enc_password",
			UserID:   1234,
		},
		{
			ID:       2,
			UID:      2,
			Website:  "https://example2.com",
			Login:    "zxc",
			Password: "enc_password2",
			UserID:   1234,
		},
	}
	repo.EXPECT().SetList(gomock.Any(), list).Return(nil)

	service := NewCredsService(repo)

	err := service.SetList(context.Background(), list)
	require.NoError(t, err)
}
