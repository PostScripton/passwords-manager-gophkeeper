package services

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"time"
)

type ISync interface {
	GetSyncAt() time.Time
}

type SyncService struct {
	client      pb.CredsClient
	credsRepo   repository.CredsSecrets
	authService Auth
}

var _ Sync = (*SyncService)(nil)

func NewSyncService(client pb.CredsClient, credsRepo repository.CredsSecrets, authService Auth) Sync {
	return &SyncService{
		client:      client,
		credsRepo:   credsRepo,
		authService: authService,
	}
}

func (s *SyncService) SyncCreds(ctx context.Context) error {
	response, err := s.client.GetAllCreds(ctx, &pb.GetAllCredsRequest{})
	if err != nil {
		return err
	}

	remoteList := s.getRemoteCreds(response.AllCreds)

	if err = s.credsRepo.SetList(ctx, remoteList); err != nil {
		return err
	}

	return nil
}

func (s *SyncService) UploadCreds(ctx context.Context) error {
	userID, err := s.authService.GetID(ctx)
	if err != nil {
		return err
	}

	localList, err := s.credsRepo.GetList(ctx, userID)
	if err != nil {
		return err
	}

	request := &pb.SetAllCredsRequest{
		AllCreds: s.getLocalCreds(localList),
	}

	if _, err = s.client.SetAllCreds(ctx, request); err != nil {
		return err
	}

	return nil
}

func (s *SyncService) getRemoteCreds(remoteList []*pb.SingleCreds) []models.CredsSecret {
	result := make([]models.CredsSecret, 0, len(remoteList))

	for _, secret := range remoteList {
		result = append(result, models.CredsSecret{
			ID:             secret.Id,
			UID:            secret.Uid,
			Website:        secret.Website,
			Login:          secret.Login,
			Password:       secret.EncPassword,
			AdditionalData: secret.AdditionalData,
			UserID:         int(secret.UserId),
		})
	}

	return result
}

func (s *SyncService) getLocalCreds(localList []*models.CredsSecret) []*pb.SingleCreds {
	result := make([]*pb.SingleCreds, 0, len(localList))

	for _, secret := range localList {
		result = append(result, &pb.SingleCreds{
			Id:             secret.ID,
			Uid:            secret.UID,
			Website:        secret.Website,
			Login:          secret.Login,
			EncPassword:    secret.Password,
			AdditionalData: secret.AdditionalData,
			UserId:         int64(secret.UserID),
		})
	}

	return result
}
