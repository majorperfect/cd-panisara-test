package tls

import (
	"crypto/tls"
	"errors"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadTLSCertForServer() (credentials.TransportCredentials, error) {
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile == "" || keyFile == "" {
		return nil, errors.New("TLS_CERT_FILE or TLS_KEY_FILE is not set")
	}

	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   0,
	}

	return credentials.NewTLS(config), nil
}
