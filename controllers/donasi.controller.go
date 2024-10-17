package controllers

import (
	"net/http"

	"github.com/alpinn/auth-go/config"
	models "github.com/alpinn/auth-go/models"
	services "github.com/alpinn/auth-go/services"
	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
)

func CreateDonasi(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var donasi models.Donasi

		if err := c.ShouldBindJSON(&donasi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		sessionKey := c.Request.Header.Get("Session-Key")
		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "No active session found"})
			return
		}

		email, err := services.Rdb.Get(config.Ctx, sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Please login again"})
			return
		}

		var user models.User
		err = db.QueryRow("SELECT id, name FROM users WHERE email = :1", email).Scan(&user.ID, &user.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch user details"})
			return
		}

		donasi.UserID = user.ID
		donasi.Name = user.Name

		err = services.PostDonasi(db, donasi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Donation created successfully"})
	}
}

func GetAllDonasi(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionKey := c.Request.Header.Get("Session-Key")
		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "No active session found"})
			return
		}

		email, err := services.Rdb.Get(config.Ctx, sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Please login again"})
			return
		}

		var userRole string
		err = db.QueryRow("SELECT role FROM users WHERE email = :1", email).Scan(&userRole)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve user role"})
			return
		}

		donasis, err := services.GetDonasi(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, donasis)
	}
}
