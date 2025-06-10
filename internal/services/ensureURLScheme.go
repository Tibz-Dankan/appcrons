package services

import (
	"strings"
)

func EnsureURLScheme(url string) string {
	if strings.HasPrefix(strings.ToLower(url), "http://") ||
		strings.HasPrefix(strings.ToLower(url), "https://") {
		return url
	}

	return "https://" + url
}
