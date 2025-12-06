package controllers

import (
	"auth-gorm/config"
	"auth-gorm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)
	c.JSON(200, users)
}

func GetById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, "user not found")
	}
	c.JSON(200, user)
}

func Create(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
	c.JSON(200, user)
}

func Update(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(404, user)
	}

	var input models.User
	c.ShouldBindJSON(&input)

	config.DB.Model(&user).Updates(input)
	c.JSON(200, gin.H{"message": "user updated", "user": user})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.User{}, id)
	c.JSON(200, gin.H{"message": "user updated"})
}
