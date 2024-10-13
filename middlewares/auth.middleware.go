package middlewares

import (
	"net/http"

	"github.com/alpinn/auth-go/config"
	"github.com/alpinn/auth-go/services"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AdminMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionKey := c.Request.Header.Get("Session-Key")
		if sessionKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "No active session found"})
			c.Abort()
			return
		}

		email, err := services.Rdb.Get(config.Ctx, sessionKey).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Please login"})
			c.Abort()
			return
		}

		var userRole string
		err = db.QueryRow("SELECT role FROM users WHERE email = $1", email).Scan(&userRole)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve user role"})
			c.Abort()
			return
		}

		if userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"msg": "Access denied."})
			c.Abort()
			return
		}

		c.Next()
	}
}
