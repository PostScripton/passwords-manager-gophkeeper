package crypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const privateKeyType = "PRIVATE KEY"

func PrivateKeyToPemBytes(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: privateKeyType, Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
}

func PrivateKeyFromPemBytes(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != privateKeyType {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed parsing private key from PEM block: %w", err)
	}

	return key, nil
}
