package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Host    string        `env:"HOST"`
	Port    string        `env:"PORT"`
	MongoDB MongoDBConfig `envPrefix:"MONGODB_"`
}

type MongoDBConfig struct {
	URI    string `env:"URI"`
	DBName string `env:"DB_NAME"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
