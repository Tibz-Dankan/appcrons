package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func SignJWTToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(3 * time.Hour).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	var jwtSecretKey = []byte(secretKey)

	accessToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return accessToken, nil
}
