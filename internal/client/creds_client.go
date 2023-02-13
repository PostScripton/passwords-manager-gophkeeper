package client

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/interceptor"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCredsClient(ctx context.Context, address string, settingsRepo repository.Settings) pb.CredsClient {
	authInterceptor := interceptor.NewUnaryClientAuthInterceptor(settingsRepo)

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Handle()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Connecting to gRPC server")
	}

	go func() {
		<-ctx.Done()
		if err = conn.Close(); err != nil {
			log.Fatal().Err(err).Msg("Closing gRPC connection")
		}
	}()

	return pb.NewCredsClient(conn)
}
