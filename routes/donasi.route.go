package routes

import (
	"database/sql"

	controllers "github.com/alpinn/auth-go/controllers"
	"github.com/gin-gonic/gin"
)

func DonasiRouter(r *gin.Engine, db *sql.DB) {
	r.GET("/me-donasi", controllers.GetAllDonasi(db))
	r.POST("/beri-donasi", controllers.CreateDonasi(db))
}
