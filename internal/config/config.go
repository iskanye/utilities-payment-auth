package config

import (
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
