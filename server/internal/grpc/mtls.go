package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// LoadServerTLSCredentials loads certificates for mTLS server
func LoadServerTLSCredentials(serverCert, serverKey, clientCA string) (grpc.ServerOption, error) {
	// Load server certificate
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load server key pair: %w", err)
	}

	// Load client CA
	ca, err := os.ReadFile(clientCA)
	if err != nil {
		return nil, fmt.Errorf("failed to read client CA: %w", err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("failed to append client CA")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    capool,
	}

	return grpc.Creds(credentials.NewTLS(config)), nil
}

// LoadClientTLSCredentials loads certificates for mTLS client
func LoadClientTLSCredentials(clientCert, clientKey, serverCA string) (grpc.DialOption, error) {
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load client key pair: %w", err)
	}

	// Load server CA
	ca, err := os.ReadFile(serverCA)
	if err != nil {
		return nil, fmt.Errorf("failed to read server CA: %w", err)
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("failed to append server CA")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      capool,
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(config)), nil
}
