package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OracleDSN string
	RedisAddr string
}

var Ctx = context.Background()

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		// PostgresDSN: os.Getenv("POSTGRES_DSN"),
		OracleDSN: os.Getenv("ORACLE_DSN"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
	}
}
