package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func SignJWTToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(9 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	var jwtSecretKey = []byte(secretKey)

	accessToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return accessToken, nil
}
