package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	r.POST("/refresh", controllers.RefreshToken)
	r.GET("/me", controllers.RequireAuth, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "login sukses"})
	})
	r.POST("/logout", controllers.RequireAuth, controllers.Logout, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "logout sukses"})
	})
}
