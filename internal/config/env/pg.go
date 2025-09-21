package env

import (
	"os"

	"github.com/based-chat/auth/internal/config"
)

var _ config.PostgresConfig = (*PostgresConfig)(nil)

const (
	envPostgresDSN = "POSTGRES_DSN"
)

type PostgresConfig struct {
	dsn string
}

func (p *PostgresConfig) DSN() string {
	return p.dsn
}

func NewPostgresConfig() (*PostgresConfig, error) {
	dsn := os.Getenv(envPostgresDSN)
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
	}

	return &PostgresConfig{
		dsn: dsn,
	}, nil
}
