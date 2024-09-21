package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Guard(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := c.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		log.Printf("Token extracted: %s", tokenString)

		token, err := verifyToken(tokenString, secret)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Println("Invalid token claims or token is not valid")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Println("Token claims:", claims)

		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			log.Println("userID not found in token claims or is not a valid type")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID := uint(userIDFloat)
		log.Printf("Token verified successfully. UserID: %v", userID)

		c.Set("user_id", userID)

		c.Next()
	}
}

func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
