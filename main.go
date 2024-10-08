package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/alpinn/auth-go/config"
	"github.com/alpinn/auth-go/routes"
	"github.com/alpinn/auth-go/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	// _ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	db  *sqlx.DB
	Rdb *redis.Client
	ctx = context.Background()
)

func main() {
	// Load configuration
	cfg := config.Load()

	// postgres connection
	var err error
	db, err = sqlx.Connect("pgx", cfg.PostgresDSN)
	log.Printf("Connecting to PostgreSQL with DSN: %s", cfg.PostgresDSN)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping Postgres:", err)
	} else {
		log.Println("Successfully connected to Postgres!")
	}

	// Redis connection
	Rdb = redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer Rdb.Close()

	// Test the Redis connection with Ping
	_, err = Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	} else {
		log.Println("Successfully connected to Redis!")
	}
	services.InitRedis(Rdb)

	// Gin setup
	r := gin.Default()

	routes.AuthRouter(r, db)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})
	log.Println("Server running on localhost:8080")
	r.Run(":8080")
}
