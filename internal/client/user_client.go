package client

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(ctx context.Context, address string) pb.UserClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("Connecting to gRPC server")
	}

	go func() {
		<-ctx.Done()
		if err = conn.Close(); err != nil {
			log.Fatal().Err(err).Msg("Closing gRPC connection")
		}
	}()

	return pb.NewUserClient(conn)
}
