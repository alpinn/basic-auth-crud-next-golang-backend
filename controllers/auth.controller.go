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
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Role            string `json:"role"`
}

func Register(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody RegisterRequest

		// Bind JSON input to requestBody struct
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		if requestBody.Password != requestBody.PasswordConfirm {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
			return
		}

		// Map data to the User model
		user := models.User{
			ID:       uuid.New(),
			Name:     requestBody.Name,
			Email:    requestBody.Email,
			Password: requestBody.Password,
			Role:     requestBody.Role,
		}

		err := services.RegisterUser(db, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
	}
}

func GetAllUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionKey := c.Request.Header.Get("Session-Key")
		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Try to login"})
			return
		}

		email, err := services.Rdb.Get(config.Ctx, sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Session expired, please login again"})
			return
		}

		var user models.User
		err = db.QueryRow("SELECT id, name, email, role FROM public.users WHERE email = $1", email).
			Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch user details"})
			}
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"msg": "Access denied"})
			return
		}

		users, err := services.GetUsers(db)
		if err != nil {
			log.Printf("GetAllUser: failed to get users: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to get users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func Login(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Password received from request: '%s'", credentials.Password)

		user, err := services.LoginUser(db, credentials.Email, credentials.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		sessionKey := user.ID.String()
		err = services.Rdb.Set(config.Ctx, sessionKey, user.Email, 30*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":       user,
			"sessionKey": sessionKey,
		})
	}
}

func Me(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the session key from the request headers
		sessionKey := c.Request.Header.Get("Session-Key")

		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Try to login"})
			return
		}

		email, err := services.Rdb.Get(config.Ctx, sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Try to login"})
			return
		}

		var user models.User
		err = db.QueryRow("SELECT id, name, email, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"msg": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal Server Error"})
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func Logout(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionKey := c.Request.Header.Get("Session-Key")

		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "No active session found"})
			return
		}

		err := services.Rdb.Del(config.Ctx, sessionKey).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to log out"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Successfully logged out"})
	}
}
