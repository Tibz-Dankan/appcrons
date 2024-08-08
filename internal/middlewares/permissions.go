package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	services "github.com/Tibz-Dankan/keep-active/internal/services"
)

type PermissionID struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Permission string `json:"permission"`
}

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
			services.AppError("You do not have permission to perform this action", 403, w)
			return
		}

		permissionID := getPermissionID(r)

		hasPermission := hasPermissions(permissionID.Permission, userPermissions.Permissions)

		if permissionID.Type == "user" && permissionID.ID != userPermissions.UserID || !hasPermission {
			services.AppError("You do not have permission to perform this action", 403, w)
			return
		}
		if permissionID.Type == "app" {
			var found bool = false
			for _, app := range userPermissions.Apps {
				if permissionID.ID == app.ID {
					found = true
					break
				}
			}
			if !found {
				services.AppError("You do not have permission to perform this action", 403, w)
				return
			}
			if !hasPermission {
				services.AppError("You do not have permission to perform this action", 403, w)
				return
			}
		}
		if permissionID.Type == "requestTime" {
			var found bool = false
			for _, app := range userPermissions.Apps {
				for _, requestTime := range app.RequestTimes {
					if permissionID.ID == requestTime.ID {
						found = true
						break
					}
				}
			}
			if !found {
				services.AppError("You do not have permission to perform this action", 403, w)
				return
			}
			if !hasPermission {
				services.AppError("You do not have permission to perform this action", 403, w)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserPermissionsKey, userPermissions)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func hasPermissions(candidatePermission string, permissions []string) bool {
	var hasPermission bool = false
	for _, permission := range permissions {
		if permission == candidatePermission {
			hasPermission = true
			break
		}
	}
	return hasPermission
}

func getPermissionID(r *http.Request) PermissionID {
	permissionID := PermissionID{}

	isAppRoute := strings.Contains(r.URL.Path, "/api/v1/apps")
	if isAppRoute {
		return getPermissionIDForAppRoutes(r)
	}

	isRequestRoute := strings.Contains(r.URL.Path, "/api/v1/requests")
	if isRequestRoute {
		return getPermissionIDForRequestRoutes(r)
	}

	isAuthRoute := strings.Contains(r.URL.Path, "/api/v1/auth")
	if isAuthRoute {
		return getPermissionIDForAuthRoutes(r)
	}

	isFeedbackRoute := strings.Contains(r.URL.Path, "/api/v1/feedback")
	if isFeedbackRoute {
		return getPermissionIDForFeedbackRoutes(r)
	}

	return permissionID
}

func getPermissionIDForAppRoutes(r *http.Request) PermissionID {
	permissionID := PermissionID{}

	getAppRegex := regexp.MustCompile(`^/api/v1/apps/get/([a-zA-Z0-9-]+)$`)
	getAppMatchString := getAppRegex.FindStringSubmatch(r.URL.Path)
	if len(getAppMatchString) > 1 {
		appId := getAppMatchString[1]
		permissionID.ID = appId
		permissionID.Type = "app"
		permissionID.Permission = "READ"
		return permissionID
	}

	postAppRegex := regexp.MustCompile(`^/api/v1/apps/post$`)
	postAppMatchString := postAppRegex.FindStringSubmatch(r.URL.Path)
	if len(postAppMatchString) > 0 {
		userId, _ := r.Context().Value(UserIDKey).(string)
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "WRITE"
		return permissionID
	}

	updateAppRegex := regexp.MustCompile(`^/api/v1/apps/update/([a-zA-Z0-9-]+)$`)
	updateAppMatchString := updateAppRegex.FindStringSubmatch(r.URL.Path)
	if len(updateAppMatchString) > 1 {
		appId := updateAppMatchString[1]
		permissionID.ID = appId
		permissionID.Type = "app"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	getAppByUserRegex := regexp.MustCompile(`^/api/v1/apps/get-by-user$`)
	getAppByUserMatchString := getAppByUserRegex.FindStringSubmatch(r.URL.Path)
	if len(getAppByUserMatchString) > 0 {
		userId := r.URL.Query().Get("userId")
		log.Println("userId:", userId)
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "READ"
		return permissionID
	}

	disableAppRegex := regexp.MustCompile(`^/api/v1/apps/disable/([a-zA-Z0-9-]+)$`)
	disableAppMatchString := disableAppRegex.FindStringSubmatch(r.URL.Path)
	if len(disableAppMatchString) > 1 {
		appId := disableAppMatchString[1]
		permissionID.ID = appId
		permissionID.Type = "app"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	enableAppRegex := regexp.MustCompile(`^/api/v1/apps/enable/([a-zA-Z0-9-]+)$`)
	enableAppMatchString := enableAppRegex.FindStringSubmatch(r.URL.Path)
	if len(enableAppMatchString) > 1 {
		appId := enableAppMatchString[1]
		permissionID.ID = appId
		permissionID.Type = "app"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	deleteAppRegex := regexp.MustCompile(`^/api/v1/apps/delete/([a-zA-Z0-9-]+)$`)
	deleteAppMatchString := deleteAppRegex.FindStringSubmatch(r.URL.Path)
	if len(deleteAppMatchString) > 1 {
		appId := deleteAppMatchString[1]
		permissionID.ID = appId
		permissionID.Type = "app"
		permissionID.Permission = "DELETE"
		return permissionID
	}

	searchAppRegex := regexp.MustCompile(`^/api/v1/apps/search$`)
	searchAppMatchString := searchAppRegex.FindStringSubmatch(r.URL.Path)
	if len(searchAppMatchString) > 0 {
		userId := r.URL.Query().Get("userId")
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "READ"
		return permissionID
	}

	getAppsLastRequestByUserRegex := regexp.MustCompile(`^/api/v1/apps/get-apps-last-request-by-user$`)
	getAppsLastRequestByUserMatchString := getAppsLastRequestByUserRegex.FindStringSubmatch(r.URL.Path)
	if len(getAppsLastRequestByUserMatchString) > 0 {
		userId := r.URL.Query().Get("userId")
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "READ"
		return permissionID
	}

	return permissionID
}

func getPermissionIDForRequestRoutes(r *http.Request) PermissionID {
	permissionID := PermissionID{}

	getRequestsByAppRegex := regexp.MustCompile(`^/api/v1/requests/get-by-app$`)
	getRequestsByAppMatchString := getRequestsByAppRegex.FindStringSubmatch(r.URL.Path)
	if len(getRequestsByAppMatchString) > 0 {
		appId := r.URL.Query().Get("appId")
		permissionID.ID = appId
		permissionID.Type = "app"
		permissionID.Permission = "READ"
		return permissionID
	}

	getLiveRequestRegex := regexp.MustCompile(`^/api/v1/requests/get-live$`)
	getLiveRequestMatchString := getLiveRequestRegex.FindStringSubmatch(r.URL.Path)
	if len(getLiveRequestMatchString) > 0 {
		userId, _ := r.Context().Value(UserIDKey).(string)
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "READ"
		return permissionID
	}

	postRequestTimeRegex := regexp.MustCompile(`^/api/v1/requests/post-request-time$`)
	postRequestTimeMatchString := postRequestTimeRegex.FindStringSubmatch(r.URL.Path)
	if len(postRequestTimeMatchString) > 0 {
		requestTime := models.RequestTime{}
		err := json.NewDecoder(r.Body).Decode(&requestTime)
		if err != nil {
			log.Println("Error decoding appId:", err)
		}
		permissionID.ID = requestTime.AppID
		permissionID.Type = "app"
		permissionID.Permission = "WRITE"
		return permissionID
	}

	updateRequestTimeRegex := regexp.MustCompile(`^/api/v1/requests/update-request-time/([a-zA-Z0-9-]+)$`)
	updateRequestTimeMatchString := updateRequestTimeRegex.FindStringSubmatch(r.URL.Path)
	if len(updateRequestTimeMatchString) > 1 {
		requestTimeId := updateRequestTimeMatchString[1]
		permissionID.ID = requestTimeId
		permissionID.Type = "requestTime"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	deleteRequestTimeRegex := regexp.MustCompile(`^/api/v1/requests/delete-request-time/([a-zA-Z0-9-]+)$`)
	deleteRequestTimeMatchString := deleteRequestTimeRegex.FindStringSubmatch(r.URL.Path)
	if len(deleteRequestTimeMatchString) > 1 {
		requestTimeId := deleteRequestTimeMatchString[1]
		permissionID.ID = requestTimeId
		permissionID.Type = "requestTime"
		permissionID.Permission = "DELETE"
		return permissionID
	}

	return permissionID
}

func getPermissionIDForAuthRoutes(r *http.Request) PermissionID {
	permissionID := PermissionID{}

	updatePasswordRegex := regexp.MustCompile(`^/api/v1/auth/user/update-password/([a-zA-Z0-9-]+)$`)
	updatePasswordMatchString := updatePasswordRegex.FindStringSubmatch(r.URL.Path)
	if len(updatePasswordMatchString) > 1 {
		userId := updatePasswordMatchString[1]
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	updateUserRegex := regexp.MustCompile(`^/api/v1/auth/user/update/([a-zA-Z0-9-]+)$`)
	updateUserMatchString := updateUserRegex.FindStringSubmatch(r.URL.Path)
	if len(updateUserMatchString) > 1 {
		userId := updateUserMatchString[1]
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	return permissionID
}

// TODO: to add user feedback in the permission cache

func getPermissionIDForFeedbackRoutes(r *http.Request) PermissionID {
	permissionID := PermissionID{}

	postFeedbackRegex := regexp.MustCompile(`^/api/v1/feedback/post$`)
	postFeedbackMatchString := postFeedbackRegex.FindStringSubmatch(r.URL.Path)
	if len(postFeedbackMatchString) > 0 {
		userId, _ := r.Context().Value(UserIDKey).(string)
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "WRITE"
		return permissionID
	}

	getFeedByUserRegex := regexp.MustCompile(`^/api/v1/feedback/get-by-user$`)
	getFeedbackByUserMatchString := getFeedByUserRegex.FindStringSubmatch(r.URL.Path)
	if len(getFeedbackByUserMatchString) > 0 {
		userId := r.URL.Query().Get("userId")
		permissionID.ID = userId
		permissionID.Type = "user"
		permissionID.Permission = "READ"
		return permissionID
	}

	updateFeedbackRegex := regexp.MustCompile(`^/api/v1/feedback/update/([a-zA-Z0-9-]+)$`)
	updateFeedbackMatchString := updateFeedbackRegex.FindStringSubmatch(r.URL.Path)
	if len(updateFeedbackMatchString) > 1 {
		feedbackId := updateFeedbackMatchString[1]
		permissionID.ID = feedbackId
		permissionID.Type = "feedback"
		permissionID.Permission = "EDIT"
		return permissionID
	}

	return permissionID
}
