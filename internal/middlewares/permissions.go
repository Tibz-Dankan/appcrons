package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	services "github.com/Tibz-Dankan/keep-active/internal/services"
)

type PermissionContextKey string

const UserPermissionsKey PermissionContextKey = "permissionId"

func HasPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(UserIDKey).(string)
		if !ok {
			log.Println("UserID not found in context")
			services.AppError("UserID not found in context", 500, w)
			return
		}

		permissions := models.Permissions{}

		userPermissions, err := permissions.Get(userId)
		if err != nil {
			log.Println("Error reading user permissions:", err)
			services.AppError("UserID not found in context", 500, w)

			return
		}

		if userPermissions.UserID == "" {
			log.Println("User has no permissions")
			services.AppError("User has no permissions to perform this action", 403, w)
			return
		}
		// TODO: To implement more robust permission checks for every authenticated request
		ctx := context.WithValue(r.Context(), UserPermissionsKey, userPermissions)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
