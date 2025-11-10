package config

import (
	"os"
	"time"
)

type Config struct {
	StoragePath    string     `yaml:"storage_path" env-required:"true"`
	Secret         string     `yaml:"secret"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func (c *Config) LoadSecret() {
	// Если секрет не задан в конфиге ищем его в параметрах окружения
	if c.Secret == "" {
		c.Secret = os.Getenv("AUTH_SECRET")
	}
	if c.Secret == "" {
		panic("auth secret mustnt be empty")
	}
}
