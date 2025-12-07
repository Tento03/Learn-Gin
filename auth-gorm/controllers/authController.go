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

/* ===================== REGISTER ======================= */
func Register(c *gin.Context) {
	var body models.Auth
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	user := models.Auth{
		Username: body.Username,
		Password: string(hashed),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "register berhasil"})
}

/* ===================== LOGIN ======================= */
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

	// ACCESS TOKEN (15 menit)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	accessString, _ := accessToken.SignedString(jwtSecret)

	// REFRESH TOKEN (7 hari)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshString, _ := refreshToken.SignedString(jwtSecret)

	// Simpan refresh token
	user.RefreshToken = refreshString
	config.DB.Save(&user)

	c.JSON(200, gin.H{
		"message":       "login berhasil",
		"access_token":  accessString,
		"refresh_token": refreshString,
	})
}

/* ===================== REFRESH TOKEN ======================= */
func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "refresh_token diperlukan"})
		return
	}

	// Parse Refresh Token
	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "refresh token invalid"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(float64))

	// cek refresh token di DB
	var user models.Auth
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(401, gin.H{"error": "user tidak ditemukan"})
		return
	}

	if user.RefreshToken != body.RefreshToken {
		c.JSON(401, gin.H{"error": "refresh token tidak cocok"})
		return
	}

	// BUAT ACCESS TOKEN BARU
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})

	accessString, _ := accessToken.SignedString(jwtSecret)

	c.JSON(200, gin.H{"access_token": accessString})
}

/* ===================== LOGOUT (hapus refresh token) ======================= */
func Logout(c *gin.Context) {
	userID := c.GetUint("userID")

	var user models.Auth
	config.DB.First(&user, userID)

	user.RefreshToken = ""
	config.DB.Save(&user)

	c.JSON(200, gin.H{"message": "logout berhasil"})
}

/* ===================== MIDDLEWARE ======================= */
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

	// Ambil UserID dari token
	claims := token.Claims.(jwt.MapClaims)
	c.Set("userID", uint(claims["id"].(float64)))

	c.Next()
}
