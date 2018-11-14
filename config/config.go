package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Listen        string `env:"LISTEN" envDefault:"localhost:8080"`
	PgSQL         string `env:"PGSQL" envDefault:"postgres://valynok:qwedfgq7@51.15.112.136:5432/valynok_db?sslmode=disable"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFileName   string `env:"LOG_FILENAME" envDefault:"main.log"`
	DirectoryPath string `env:"DIRECTORY_PATH" envDefault:"./files"`
	Init          bool   `env:"INIT" envDefault:"false"`
}

func Load() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
