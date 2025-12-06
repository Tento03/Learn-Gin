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
	var user models.Auth

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	input := models.Auth{
		Username: user.Username,
		Password: string(hashed),
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username sudah digunakan"})
	}

	c.JSON(200, gin.H{"message": "Registrasi", "user": input})
}

func Login(c *gin.Context) {
	var body models.Auth
	c.Bind(&body)

	var user models.Auth
	if err := config.DB.First(&user, user.Username).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(200, gin.H{"message": "Login Berhasil", "token": tokenString})
}

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(401, gin.H{"error": "token tidak ada"})
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
