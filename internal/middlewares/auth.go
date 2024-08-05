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

type ContextKey string

const UserIDKey ContextKey = "userId"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		var bearerToken string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
			headerParts := strings.SplitN(authHeader, " ", 2)
			bearerToken = headerParts[1]
		}
		if bearerToken == "" {
			services.AppError("You are not logged in! Please to gain access.", 401, w)
			return
		}

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
		claims, validJWTClaim := token.Claims.(jwt.MapClaims)
		if !validJWTClaim || !token.Valid {
			services.AppError("Invalid Token. please login again", 403, w)
			return
		}

		if userIDClaim, ok := claims["userId"].(string); ok {
			userId = userIDClaim
		}

		User := models.User{}

		user, err := User.FindOne(userId)
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
