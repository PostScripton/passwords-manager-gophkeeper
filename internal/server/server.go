package server

import (
	"crypto/tls"
	pb "github.com/PostScripton/passwords-manager-gophkeeper/api/proto"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/interceptor"
	servicesPkg "github.com/PostScripton/passwords-manager-gophkeeper/internal/services"
	"github.com/PostScripton/passwords-manager-gophkeeper/pkg/cert"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	listener net.Listener
	core     *grpc.Server
}

func NewServer(address string, services *servicesPkg.Services, sslCertPath, sslKeyPath string) *Server {
	tlsConf, err := cert.LoadServerCertificate(sslCertPath, sslKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Loading server TLS cert")
	}

	listen, err := tls.Listen("tcp", address, tlsConf)
	if err != nil {
		log.Fatal().Err(err).Str("address", address).Msg("Listening to TCP address")
	}

	authInterceptor := interceptor.NewUnaryServerAuthInterceptor(services.Auth)

	core := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Handle()))
	pb.RegisterUserServer(core, &UserServer{
		services: services,
	})
	pb.RegisterCredsServer(core, &CredsServer{
		services: services,
	})

	return &Server{
		listener: listen,
		core:     core,
	}
}

func (s *Server) Run() error {
	return s.core.Serve(s.listener)
}

func (s *Server) Shutdown() error {
	s.core.GracefulStop()

	return nil
}
