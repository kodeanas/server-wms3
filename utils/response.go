package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse defines standard API response.
type APIResponse struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorItem for validation errors.
type ErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// SendSuccess sends a standard success response.
func SendSuccess(c *gin.Context, data interface{}, message string, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	if message == "" {
		message = "success"
	}

	c.JSON(code, APIResponse{
		Code:    code,
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendPaginatedSuccess sends paginated success response.
func SendPaginatedSuccess(c *gin.Context, data interface{}, page, limit, totalItems, totalPages int64, message string) {
	if message == "" {
		message = "success"
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: message,
		Data:    data,
		Meta: map[string]interface{}{
			"pagination": map[string]interface{}{
				"page":        page,
				"limit":       limit,
				"total_items": totalItems,
				"total_pages": totalPages,
			},
		},
	})
}

// SendError sends an error response.
func SendError(c *gin.Context, code int, message string) {
	if code == 0 {
		code = http.StatusBadRequest
	}
	if message == "" {
		message = "error"
	}

	c.JSON(code, APIResponse{
		Code:    code,
		Success: false,
		Message: message,
		Data:    nil,
		Meta:    nil,
	})
}

// SendValidationError sends validation error response.
func SendValidationError(c *gin.Context, errors []ErrorItem) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Code:    http.StatusBadRequest,
		Success: false,
		Message: "Validasi gagal",
		Data:    nil,
		Meta: map[string]interface{}{
			"errors": errors,
		},
	})
}
