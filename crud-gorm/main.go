package main

import (
	"crud-gorm/config"
	"crud-gorm/models"
	"crud-gorm/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// Auto migrate models
	config.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	routes.UserRoute(r)

	r.Run(":8080")
}
