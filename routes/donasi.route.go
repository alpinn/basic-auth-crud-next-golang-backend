package routes

import (
	controllers "github.com/alpinn/auth-go/controllers"
	"github.com/alpinn/auth-go/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DonasiRouter(r *gin.Engine, db *sqlx.DB) {
	r.GET("/donasi", middlewares.AdminMiddleware(db), controllers.GetAllDonasi(db))
	r.POST("/beri-donasi", controllers.CreateDonasi(db))
}
