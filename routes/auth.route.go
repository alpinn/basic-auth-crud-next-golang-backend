package routes

import (
	"database/sql"

	controllers "github.com/alpinn/auth-go/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine, db *sql.DB) {
	r.POST("/register", controllers.Register(db))
	r.POST("/login", controllers.Login(db))
	r.GET("/me", controllers.Me(db))
}
