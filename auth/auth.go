package auth

import (
	"api-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

const (
	USER     = "admin"
	PASSWORD = "1234"
	SECRET   = "secret"
)

func LoginHandler(c *gin.Context) {
	var user models.Credential
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	if user.Username != USER || user.Password != PASSWORD {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email atau Password salah",
		})
		return
	}

	claim := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		Issuer:    "test",
		IssuedAt:  time.Now().Unix(),
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := sign.SignedString([]byte(SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}
