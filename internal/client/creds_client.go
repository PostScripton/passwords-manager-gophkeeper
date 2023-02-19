package client

import (
	"context"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/interceptor"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"github.com/PostScripton/passwords-manager-gophkeeper/pkg/cert"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func NewCredsClient(
	ctx context.Context,
	address string,
	settingsRepo repository.Settings,
	sslCertPath, sslKeyPath string,
) pb.CredsClient {
	tlsCredential, err := cert.LoadClientCertificate(sslCertPath, sslKeyPath)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("cert-path", sslCertPath).
			Str("key-path", sslKeyPath).
			Msg("Loading client TLS cert")
	}

	authInterceptor := interceptor.NewUnaryClientAuthInterceptor(settingsRepo)

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(tlsCredential),
		grpc.WithUnaryInterceptor(authInterceptor.Handle()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Connecting to gRPC server")
	}

	go func() {
		<-ctx.Done()
		if err = conn.Close(); err != nil {
			log.Fatal().Err(err).Str("context-error", ctx.Err().Error()).Msg("Closing gRPC connection")
		}
	}()

	return pb.NewCredsClient(conn)
}
