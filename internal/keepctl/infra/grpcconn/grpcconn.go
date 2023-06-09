// Package grpcconn implements handy wrap around gRPC client connection.
package grpcconn

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const _DefaultTLSVersion = tls.VersionTLS13

var ErrAppendRootCert = errors.New("failed to append root certificate")

// Connection wraps gRPC client connection.
type Connection struct {
	conn *grpc.ClientConn
}

// New create and initializes new gRPC Connection.
// Warning: Since gRPC connections in Go are lazy, actual link is not established
// until first request.
func New(address, caPath string) (*Connection, error) {
	rootCA, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(rootCA) {
		return nil, ErrAppendRootCert
	}

	config := &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         _DefaultTLSVersion,
		RootCAs:            cp,
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)
	if err != nil {
		return nil, fmt.Errorf("grpcconn - New - grpc.Dial: %w", err)
	}

	return &Connection{conn}, nil
}

// Instance grants access to the underlying gRPC client connection.
func (c *Connection) Instance() *grpc.ClientConn {
	return c.conn
}

// Close stops gRPC client connection gracefully.
func (c *Connection) Close() error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("grpcconn - Close - c.conn.Close: %w", err)
	}

	return nil
}
