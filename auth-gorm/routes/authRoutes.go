package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/refresh", controllers.RefreshToken)

	auth := r.Group("/")
	auth.Use(controllers.RequireAuth)

	auth.GET("/me", controllers.RequireAuth)
	auth.GET("/logout", controllers.Logout)
}
