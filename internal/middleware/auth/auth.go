package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"userID"`
	jwt.RegisteredClaims
}

func CreateToken(username string, userID uint, secret string) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := CustomClaims{
		Username: username,
		UserID:   userID, // Use uint here
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing token")
		return "", err
	}

	return signedToken, nil
}
