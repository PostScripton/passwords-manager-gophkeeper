package client

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/pkg/cert"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func NewUserClient(ctx context.Context, address string, sslCertPath, sslKeyPath string) pb.UserClient {
	tlsCredential, err := cert.LoadClientCertificate(sslCertPath, sslKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Loading client TLS cert")
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tlsCredential))
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
