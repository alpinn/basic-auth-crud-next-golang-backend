package routes

import (
	controllers "github.com/alpinn/auth-go/controllers"
	"github.com/alpinn/auth-go/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AuthRouter(r *gin.Engine, db *sqlx.DB) {
	r.GET("/users", middlewares.AdminMiddleware(db), controllers.GetAllUser(db))

	r.POST("/register", controllers.Register(db))
	r.POST("/login", controllers.Login(db))

	r.GET("/me", controllers.Me(db))
	r.DELETE("/logout", controllers.Logout(db))
}
