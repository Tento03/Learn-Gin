package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/profile", controllers.RequireAuth, func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "anda berhasil login"})
	})
}
