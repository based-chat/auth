package config

import "github.com/joho/godotenv"

func Load(path string) error {
	return godotenv.Load(path)
}

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}
