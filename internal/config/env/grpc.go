package env

import (
	"net"
	"os"

	"github.com/based-chat/auth/internal/config"
)

var _ config.GRPCConfig = (*GrpcConfig)(nil)

type GrpcConfig struct {
	host string
	port string
}

// Address возвращает адрес gRPC-сервера в формате host:port.
func (g *GrpcConfig) Address() string {
	return net.JoinHostPort(g.host, g.port)
}

// Возвращает указатель на GrpcConfig и ошибку (в текущей реализации всегда nil).
func NewGRPCConfig() (*GrpcConfig, error) {
	host := os.Getenv("GRPC_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	return &GrpcConfig{
		host: host,
		port: port,
	}, nil
}
