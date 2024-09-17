package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, secret string) (string, error) {
	// Set token to expire in 1 month (30 days)
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	// Create the JWT token with the expiration time
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{username},         // The username
		ExpiresAt: jwt.NewNumericDate(expirationTime), // Token expiration time set to 1 month
	})

	// Sign the token with the provided secret
	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing token")
		return "", err
	}

	return signedToken, nil
}
