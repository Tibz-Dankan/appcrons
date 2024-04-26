package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	services "github.com/Tibz-Dankan/keep-active/internal/services"

	"github.com/golang-jwt/jwt"
)

// ContextKey is a custom type for context keys
type ContextKey string

// UserIDKey is the key to access the user ID stored in context
const UserIDKey ContextKey = "userId"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			services.AppError("Missing Authorization header", 401, w)
			return
		}
		headerParts := strings.SplitN(authorizationHeader, " ", 2)

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			services.AppError("Invalid Authorization header", 401, w)
			return
		}
		bearerToken := headerParts[1]
		secretKey := os.Getenv("JWT_SECRET")
		var jwtSecretKey = []byte(secretKey)

		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				services.AppError("Unexpected signing method", 403, w)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil {
			services.AppError(err.Error(), 403, w)
			return
		}

		var userId string
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userIDClaim, _ := claims["userId"].(string)
			userId = userIDClaim
		} else {
			services.AppError("Invalid Token, please login again", 403, w)
			return
		}

		User := models.User{}

		user, err := User.FindOne(userId)
		// TODO: To add redis read functionality here
		if err != nil {
			services.AppError(err.Error(), 500, w)
			return
		}

		if user.ID == "" {
			services.AppError("The user belonging to this token no longer exist!", 403, w)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
