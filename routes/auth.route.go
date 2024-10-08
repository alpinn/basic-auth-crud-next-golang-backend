package routes

import (
	controllers "github.com/alpinn/auth-go/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AuthRouter(r *gin.Engine, db *sqlx.DB) {
	r.POST("/register", controllers.Register(db))
	r.POST("/login", controllers.Login(db))
	r.GET("/me", controllers.Me(db))
	r.GET("/users", controllers.GetAllUser(db))
}
