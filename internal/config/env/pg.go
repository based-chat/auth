package env

import (
	"log"
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

// DSN возвращает строку подключения к PostgreSQL.
func (p *PostgresConfig) DSN() string {
	return p.dsn
}

// NewPostgresConfig создаёт и возвращает конфигурацию PostgreSQL.
// DSN читается из переменной окружения POSTGRES_DSN; если переменная не задана или пустая,
// выводится предупреждение в лог.
// Возвращает указатель на PostgresConfig и ошибку (в текущей реализации всегда nil).
func NewPostgresConfig() (*PostgresConfig, error) {
	dsn := os.Getenv(envPostgresDSN)
	if dsn == "" {
		log.Default().Println("POSTGRES_DSN is not set")
	}

	return &PostgresConfig{
		dsn: dsn,
	}, nil
}
