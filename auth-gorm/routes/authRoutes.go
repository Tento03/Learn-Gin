package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshToken)
	r.GET("/me", controllers.RequireAuth, func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "login sukses"})
	})
	r.POST("/logout", controllers.RequireAuth, func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "logout"})
	})
}
