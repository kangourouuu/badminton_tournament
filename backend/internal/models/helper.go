package models

import (
	"strings"
	"regexp"
)

func SanitizeNameForEmail(name string) string {
	// Lowercase
	s := strings.ToLower(name)
	// Remove accents (simplistic approach, or just strip non-alphanumeric)
	reg, _ := regexp.Compile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "")
	return s
}
