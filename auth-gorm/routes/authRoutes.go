package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/")
	auth.Use(controllers.RequireAuth)
	auth.GET("/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "kamu berhasil login"})
	})
}
