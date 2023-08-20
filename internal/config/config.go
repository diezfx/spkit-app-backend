package config

import (
	"os"

	"github.com/diezfx/split-app-backend/pkg/sqlite"
)

type Environment string

const (
	LocalEnv       Environment = "local"
	DevelopmentEnv Environment = "dev"
)

type Config struct {
	Addr        string
	Environment Environment
	LogLevel    string

	DB sqlite.Config
}

func Load() Config {

	env := os.Getenv("ENVIRONMENT")
	if env == string(DevelopmentEnv) {
		return Config{
			Addr:        "localhost:80",
			Environment: DevelopmentEnv,
			LogLevel:    "debug",
			DB:          sqlite.Config{Path: "ent.db", InMemory: false},
		}
	}

	return Config{
		Addr:        "localhost:5002",
		Environment: LocalEnv,
		LogLevel:    "debug",
		DB:          sqlite.Config{Path: "ent.db", InMemory: false},
	}
}
