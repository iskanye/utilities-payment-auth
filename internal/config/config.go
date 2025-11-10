package config

import (
	"os"
	"time"
)

type Config struct {
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	Secret         string
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func (c *Config) MustGetSecret() {
	// Ищем секрет
	c.Secret = os.Getenv("AUTH_SECRET")
	if c.Secret == "" {
		panic("auth secret mustnt be empty")
	}
}
