package services

import "strings"

func ContainsZ(timeStr string) bool {
	return strings.Contains(timeStr, "Z")
}
