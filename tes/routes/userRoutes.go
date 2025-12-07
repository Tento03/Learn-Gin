package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.GET("/me", controllers.RequireAuth, func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "login berhasil"})
	})
}
