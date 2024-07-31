package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)



func generateToken(user User) (string, error) {
	godotenv.Load()

	var expirationTime time.Duration

	if user.Expires_in_seconds == 0 {
		expirationTime = 24 * time.Hour
	} else {
		expirationTime = time.Duration(user.Expires_in_seconds) * time.Second
		if expirationTime > 24*time.Hour {
			expirationTime = 24 * time.Hour
		}
	}

	myToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expirationTime)),
		Subject:   fmt.Sprintf("%d", user.ID),
	})

	tokenString, err := myToken.SignedString([]byte(Cfg.JWTSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	return tokenString, nil
}