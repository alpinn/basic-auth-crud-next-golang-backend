package config

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/godror/godror"
	"github.com/joho/godotenv"
)

var DB *sql.DB
var Ctx = context.Background()

func Init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func InitDB() {
	var err error
	dsn := os.Getenv("ORACLE_DSN")
	DB, err = sql.Open("godror", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("DB Ping failed: ", err)
	}
}

func InitRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	return redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}
