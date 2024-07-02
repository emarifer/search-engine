package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	Id                   string `json:"id"`
	User                 string `json:"user"`
	Admin                bool   `json:"role"`
	jwt.RegisteredClaims `json:"claims"`
}

func CreateNewAuthToken(id, user string, isAdmin bool) (string, error) {
	claims := AuthClaims{
		Id:    id,
		User:  user,
		Admin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			Issuer:    "https://github.com/emarifer",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	secretKey, exists := os.LookupEnv("SECRET_KEY")
	if !exists {
		log.Fatalln("🔥 SECRET_KEY cannot be found in .env file")
	}

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing the token: %s", err)
	}

	return signedToken, nil
}
