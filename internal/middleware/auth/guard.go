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
		// Extract token from cookie
		auth, err := c.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie:", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Optionally remove "Bearer " prefix
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		log.Printf("Token extracted: %s", tokenString)

		// Verify the token
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

		// Log all claims
		log.Println("Token claims:", claims)

		// Extract and convert userID from token claims
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			log.Println("userID not found in token claims or is not a valid type")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Convert userID to uint (as you prefer userID to be uint)
		userID := uint(userIDFloat)
		log.Printf("Token verified successfully. UserID: %v", userID)

		// Set userID in the context as uint
		c.Set("user_id", userID)

		c.Next()
	}
}

func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Return the verified token
	return token, nil
}
