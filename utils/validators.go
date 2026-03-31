package utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// BindJSONOrFail binds JSON and fails if error
func BindJSONOrFail(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		SendError(c, 400, "Invalid request body: "+err.Error())
		return false
	}
	return true
}

// ParseInt converts string to int
func ParseInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	return v, err
}

// GetPaginationParams extracts pagination parameters from query
func GetPaginationParams(c *gin.Context, defaultLimit int) (int, int) {
	limit := defaultLimit
	offset := 0

	if l := c.Query("limit"); l != "" {
		if v, err := ParseInt(l); err == nil && v > 0 {
			limit = v
		}
	}

	if o := c.Query("offset"); o != "" {
		if v, err := ParseInt(o); err == nil && v >= 0 {
			offset = v
		}
	}

	return limit, offset
}

// ValidateEmail checks if email is valid
func ValidateEmail(email string) bool {
	// Simple email validation
	if len(email) < 5 {
		return false
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

// ValidatePhone checks if phone is valid
func ValidatePhone(phone string) bool {
	if len(phone) < 10 {
		return false
	}
	// Remove non-digit characters
	cleaned := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		} else if char == '+' && len(cleaned) == 0 {
			cleaned += string(char)
		}
	}
	return len(cleaned) >= 10
}

// ValidateRequired checks if string is not empty
func ValidateRequired(value string) bool {
	return strings.TrimSpace(value) != ""
}

// ValidateMinLength checks if string meets minimum length
func ValidateMinLength(value string, minLength int) bool {
	return len(strings.TrimSpace(value)) >= minLength
}

// ValidateMaxLength checks if string is within maximum length
func ValidateMaxLength(value string, maxLength int) bool {
	return len(value) <= maxLength
}

// TrimSpace removes leading and trailing whitespace
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}
