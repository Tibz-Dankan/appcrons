package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Tibz-Dankan/keep-active/internal/models"

	"github.com/golang-jwt/jwt"
)

// OptionalAuth resolves userId from a bearer token if present and valid,
// but unlike Auth it never rejects the request - anonymous callers proceed
// with UserIDKey set to "".
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := ""

		authHeader := r.Header.Get("Authorization")
		var bearerToken string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
			headerParts := strings.SplitN(authHeader, " ", 2)
			if len(headerParts) > 1 {
				bearerToken = headerParts[1]
			}
		}

		if bearerToken != "" {
			secretKey := os.Getenv("JWT_SECRET")
			jwtSecretKey := []byte(secretKey)

			token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
				return jwtSecretKey, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userIDClaim, ok := claims["userId"].(string); ok && userIDClaim != "" {
						User := models.User{}
						if user, err := User.FindOne(userIDClaim); err == nil && user.ID != "" {
							userId = user.ID
						}
					}
				}
			}
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
