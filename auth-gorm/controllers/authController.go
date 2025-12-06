package controllers

import (
	"auth-gorm/config"
	"auth-gorm/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(c *gin.Context) {
	var body models.Auth

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	user := models.Auth{
		Username: body.Username,
		Password: string(hashed),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "register berhasil", "user": user})
}

func Login(c *gin.Context) {
	var body models.Auth
	c.ShouldBindJSON(&body)

	var user models.Auth
	if err := config.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	var tokenString, _ = token.SignedString(jwtSecret)
	c.JSON(200, gin.H{"message": "login berhasil", "token": tokenString})
}

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "token tidak ditemukan"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "token invalid"})
		c.Abort()
		return
	}
	c.Next()
}
