package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/alpinn/auth-go/config"
	"github.com/alpinn/auth-go/routes"
	"github.com/alpinn/auth-go/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/godror/godror"
)

var (
	db  *sql.DB
	Rdb *redis.Client
	ctx = context.Background()
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Oracle connection
	var err error
	db, err = sql.Open("godror", cfg.OracleDSN)
	if err != nil {
		log.Fatal("Failed to connect to Oracle:", err)
	}
	defer db.Close()

	// Redis connection
	Rdb = redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer Rdb.Close()

	services.InitRedis(Rdb)

	// Gin setup
	r := gin.Default()

	routes.AuthRouter(r, db)
	routes.DonasiRouter(r, db)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})
	log.Println("Server running on localhost:8080")
	r.Run(":8080")
}
