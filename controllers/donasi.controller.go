package controllers

import (
	"database/sql"
	"net/http"

	models "github.com/alpinn/auth-go/models"
	services "github.com/alpinn/auth-go/services"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func CreateDonasi(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var donasi models.Donasi
		if err := c.ShouldBindJSON(&donasi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
			return
		}

		// Parse user_id from path params
		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid user ID"})
			return
		}
		donasi.UserID = userID // Set UserID for the donasi

		err = services.PostDonasi(db, donasi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create product"})
			return
		}
		c.JSON(http.StatusOK, donasi)
	}
}

func GetAllDonasi(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		donasi, err := services.GetDonasi(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve donasi"})
			return
		}
		c.JSON(http.StatusOK, donasi)
	}
}
