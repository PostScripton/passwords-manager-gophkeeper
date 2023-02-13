package server

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/services"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

type CredsServer struct {
	services *services.Services
	pb.UnimplementedCredsServer
}

var _ pb.CredsServer = (*CredsServer)(nil)

func (c *CredsServer) GetAllCreds(ctx context.Context, _ *pb.GetAllCredsRequest) (*pb.GetAllCredsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to get user ID from context")
	}
	userID, err := strconv.Atoi(userIDs[0])
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error on converting user ID to integer")
	}

	list, err := c.services.Creds.GetList(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Getting credentials list for the user")
		return nil, status.Errorf(codes.Unknown, "Getting credentials list for the user")
	}

	out := make([]*pb.SingleCreds, 0, len(list))

	for _, secret := range list {
		out = append(out, &pb.SingleCreds{
			Id:             secret.ID,
			Uid:            secret.UID,
			Website:        secret.Website,
			Login:          secret.Login,
			EncPassword:    secret.Password,
			AdditionalData: secret.AdditionalData,
			UserId:         int64(secret.UserID),
		})
	}

	response := &pb.GetAllCredsResponse{AllCreds: out}

	return response, nil
}

func (c *CredsServer) SetAllCreds(
	ctx context.Context,
	in *pb.SetAllCredsRequest,
) (*pb.SetAllCredsResponse, error) {
	list := make([]models.CredsSecret, 0, len(in.AllCreds))

	for _, secret := range in.AllCreds {
		list = append(list, models.CredsSecret{
			ID:             secret.Id,
			UID:            secret.Uid,
			Website:        secret.Website,
			Login:          secret.Login,
			Password:       secret.EncPassword,
			AdditionalData: secret.AdditionalData,
			UserID:         int(secret.UserId),
		})
	}

	if err := c.services.Creds.SetList(ctx, list); err != nil {
		log.Error().Err(err).Msg("Setting credentials list for the user")
		return nil, status.Errorf(codes.Unknown, "Setting credentials list for the user")
	}

	return &pb.SetAllCredsResponse{}, nil
}
