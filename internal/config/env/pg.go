package env

import (
	"errors"
	"log"
	"os"

	"github.com/based-chat/auth/internal/config"
)

var _ config.PostgresConfig = (*PostgresConfig)(nil)

const (
	envPostgresDSN = "POSTGRES_DSN"
)

var (
	errPostgresDSNNotSet = errors.New("postgres dsn is not set")
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
		log.Default().Println("postgres dsn is not set")
		return nil, errPostgresDSNNotSet
	}

	return &PostgresConfig{
		dsn: dsn,
	}, nil
}
