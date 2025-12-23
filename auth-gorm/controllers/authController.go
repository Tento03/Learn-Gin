package controllers

import (
	"auth-gorm/config"
	"auth-gorm/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Auth
	if err := config.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "username gada"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(400, gin.H{"error": "password salah"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Minute),
	})
	accessString, _ := accessToken.SignedString(jwtSecret)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Hour),
	})
	refreshString, _ := refreshToken.SignedString(jwtSecret)

	user.RefreshToken = refreshString
	config.DB.Save(&user)
	c.JSON(200, gin.H{"access token": accessString, "refresh token": refreshString})
}

func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if !token.Valid || err != nil {
		c.JSON(401, gin.H{"error": "refresh token invalid"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	var user models.Auth
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(401, gin.H{"message": "user ga ditemukan"})
		return
	}

	if body.RefreshToken != user.RefreshToken {
		c.JSON(401, gin.H{"message": "user ga ditemukan"})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Minute),
	})
	accessString, _ := accessToken.SignedString(jwtSecret)
	c.JSON(200, gin.H{"access token": accessString})
}

func Logout(c *gin.Context) {
	userId := c.GetUint("userId")

	var user models.Auth
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(404, gin.H{"error": "user gada"})
		return
	}

	user.RefreshToken = ""
	config.DB.Save(&user)
}

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(401, gin.H{"err": "token kosong"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if !token.Valid || err != nil {
		c.JSON(401, gin.H{"err": "token invalid"})
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))
	c.Set("userId", userId)
	c.Next()
}
