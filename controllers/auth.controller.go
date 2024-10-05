package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/alpinn/auth-go/config"
	models "github.com/alpinn/auth-go/models"
	services "github.com/alpinn/auth-go/services"
	"github.com/gin-gonic/gin"
)

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Password != c.PostForm("password_confirm") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
			return
		}

		err := services.RegisterUser(db, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
	}
}
func GetAllUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := services.GetUsers(db)
		if err != nil {
			log.Printf("GetAllUser: failed to get user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to get user"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func Login(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials models.User
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := services.LoginUser(db, credentials.Email, credentials.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		sessionKey := "session:" + user.ID.String()
		err = services.Rdb.Set(config.Ctx, sessionKey, user.Email, 30*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":       user,
			"sessionKey": sessionKey, // Return the session key
		})
	}
}

func Me(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the session key from the request headers
		sessionKey := c.Request.Header.Get("User-ID") // Use the key you stored it under

		// Check if the session key is provided
		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Try to login"})
			return
		}

		// Retrieve the email from Redis based on the session key
		email, err := services.Rdb.Get(config.Ctx, sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Try to login"})
			return
		}

		// Fetch the user from the database using the email
		var user models.User
		err = db.QueryRow("SELECT id, name, email FROM users WHERE email = :1", email).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Server Error"})
			}
			return
		}

		// Return the user data
		c.JSON(http.StatusOK, user)
	}
}
