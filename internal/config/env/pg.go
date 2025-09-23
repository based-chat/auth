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

// NewPostgresConfig создаёт и возвращает конфигурацию PostgreSQL.
// DSN читается из переменной окружения POSTGRES_DSN; если переменная не задана или пустая,
// используется значение по умолчанию "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable".
// Возвращает указатель на PostgresConfig и ошибку (в текущей реализации всегда nil).
func NewPostgresConfig() (*PostgresConfig, error) {
	dsn := os.Getenv(envPostgresDSN)
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"
	}

	return &PostgresConfig{
		dsn: dsn,
	}, nil
}
