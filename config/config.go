package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Listen string `env:"LISTEN" envDefault:"localhost:8080"`
	//PgSQL    string `env:"PGSQL" envDefault:"postgres://postgres:111111@localhost:5432/blog?sslmode=disable"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"info"`
	LogFileName   string `env:"MY_LOG_FILENAME" envDefault:"myproject.log"`
	DirectoryPath string `env:"MY_DIRECTORY_PATH" envDefault:"./files"`
}

func Load() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
