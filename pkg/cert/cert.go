package cert

import (
	"crypto/tls"
	"google.golang.org/grpc/credentials"
)

// LoadClientCertificate returns client credential TLS by paths.
func LoadClientCertificate(sslCertPath, sslKeyPath string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(sslCertPath, sslKeyPath)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}), nil
}

// LoadServerCertificate returns server tls config by paths.
func LoadServerCertificate(sslCertPath, sslKeyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(sslCertPath, sslKeyPath)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}
