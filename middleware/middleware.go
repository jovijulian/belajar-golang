package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	_ "time"
)

const (
	SECRET = "secret"
)

func AuthValid(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Token tidak ada",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, invalid := token.Method.(*jwt.SigningMethodHMAC); !invalid {
			return nil, fmt.Errorf("Invalid token, alg: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   err.Error(),
				"message": "Token tidak valid",
			})
			c.Abort()
			return
		}

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Token kadaluarsa",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	if token.Valid {
		fmt.Println("Token Verified")
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Token tidak valid",
			"message": "Unauthorized",
		})
		c.Abort()
	}
}
