package routes

import (
	"auth-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controllers.GetAll)
	r.GET("/users/:id", controllers.GetById)
	r.POST("/users", controllers.Create)
	r.PUT("/users/:id", controllers.Update)
	r.DELETE("/users/:id", controllers.Delete)
}
