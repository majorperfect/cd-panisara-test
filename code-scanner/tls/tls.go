package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCertFile := os.Getenv("TLS_CERT_FILE")
	CertFilePath := os.Getenv("TLS_CLIENT_CERT_FILE")
	KeyCertFilePath := os.Getenv("TLS_CLIENT_CERT_KEY_FILE")

	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile(serverCertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(CertFilePath, KeyCertFilePath)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}
