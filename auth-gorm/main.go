package main

import (
	"auth-gorm/config"
	"auth-gorm/models"
	"auth-gorm/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	config.DB.AutoMigrate(models.Auth{})

	r := gin.Default()
	routes.AuthRoutes(r)
	r.Run(":8080")
}
