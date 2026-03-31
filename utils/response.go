package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Message string      `json:"message"`
}

// ErrorResponse represents error API response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SendSuccess returns success response
func SendSuccess(c *gin.Context, data interface{}, message string, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	if message == "" {
		message = "success"
	}

	c.JSON(code, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendPaginatedSuccess returns paginated success response
func SendPaginatedSuccess(c *gin.Context, data interface{}, total int64, message string) {
	if message == "" {
		message = "success"
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Total:   total,
	})
}

// SendError returns error response
func SendError(c *gin.Context, code int, message string) {
	if code == 0 {
		code = http.StatusBadRequest
	}
	if message == "" {
		message = "error"
	}

	c.JSON(code, ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	})
}

// SendNotFound returns 404 response
func SendNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	SendError(c, http.StatusNotFound, message)
}

// SendUnauthorized returns 401 response
func SendUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	SendError(c, http.StatusUnauthorized, message)
}

// SendForbidden returns 403 response
func SendForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	SendError(c, http.StatusForbidden, message)
}

// SendInternalError returns 500 response
func SendInternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	SendError(c, http.StatusInternalServerError, message)
}

// SendValidationError returns 422.
func SendValidationError(c *gin.Context, errors interface{}) {
	c.JSON(422, gin.H{"status": "fail", "errors": errors})
}
