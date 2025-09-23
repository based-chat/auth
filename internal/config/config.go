package config

import "github.com/joho/godotenv"

// Load загружает переменные окружения из файла, указанного в path.
// path — путь к файлу с переменными окружения (обычно ".env"); возвращает ошибку, полученную при попытке загрузки.
func Load(path string) error {
	return godotenv.Load(path)
}

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}
