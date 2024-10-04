package main

import (
	"log"

	"github.com/alpinn/auth-go/config"
	"github.com/alpinn/auth-go/routes"
	services "github.com/alpinn/auth-go/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	config.InitDB()
	defer func() {
		if err := config.DB.Close(); err != nil {
			log.Fatalf("Failed to close DB connection: %v", err)
		}
	}()

	redisClient := config.InitRedis()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalf("Failed to close Redis connection: %v", err)
		}
	}()

	services.InitRedis(redisClient)

	r := gin.Default()

	// Register routes
	routes.AuthRouter(r, config.DB)
	routes.DonasiRouter(r, config.DB)

	log.Println("Server running on port 8080")
	r.Run(":8080")
}
