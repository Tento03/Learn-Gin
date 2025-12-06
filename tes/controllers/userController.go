package controllers

import (
	"auth-gorm/config"
	"auth-gorm/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(c *gin.Context) {
	var body models.Auth
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	user := models.Auth{
		Username: body.Username,
		Password: string(hashed),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username sudah digunakan"})
	}
	c.JSON(200, user)
}

func Login(c *gin.Context) {
	var body models.Auth
	c.ShouldBindJSON(&body)

	var user models.Auth
	if err := config.DB.First(&user, body.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username tidak ada"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(body.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)
	c.JSON(200, gin.H{"message": "login berhasil", "token": tokenString})
}

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"message": "token tidak ada"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"message": "token invalid"})
		c.Abort()
		return
	}
	c.Next()
}
