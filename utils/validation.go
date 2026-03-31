package utils

import (
	"fmt"
	"strings"
)

// MustBePositive ensures numeric value is positive.
func MustBePositive(value int64) error {
	if value <= 0 {
		return fmt.Errorf("value must be positive")
	}
	return nil
}

// CheckStatus validates stock status value.
func CheckStatus(status string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case "good", "damaged", "unavailable", "in-transit", "on-progress", "done":
		return nil
	default:
		return fmt.Errorf("status must be one of: good, damaged, unavailable, in-transit, on-progress, done")
	}
}

// ValidateOrderStatus validates order status value.
func ValidateOrderStatus(status string) error {
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case "pending", "processing", "shipped", "delivered", "cancelled":
		return nil
	default:
		return fmt.Errorf("status must be one of: pending, processing, shipped, delivered, cancelled")
	}
}

// ValidateOrderType validates order type value.
func ValidateOrderType(orderType string) error {
	orderType = strings.ToLower(strings.TrimSpace(orderType))
	switch orderType {
	case "regular", "wcargo", "xpedx", "same-day":
		return nil
	default:
		return fmt.Errorf("order type must be one of: regular, wcargo, xpedx, same-day")
	}
}
