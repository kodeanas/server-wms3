package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug creates a slug from a string
func GenerateSlug(input string) string {
	// Lowercase
	slug := strings.ToLower(input)
	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove non-alphanumeric and non-hyphen characters
	re := regexp.MustCompile(`[^a-z0-9-]`)
	slug = re.ReplaceAllString(slug, "")
	// Remove multiple hyphens
	re2 := regexp.MustCompile(`-+`)
	slug = re2.ReplaceAllString(slug, "-")
	// Trim hyphens
	slug = strings.Trim(slug, "-")
	return slug
}
