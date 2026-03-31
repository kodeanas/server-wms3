package utils

import (
	"errors"
	"fmt"
)

// CustomError represents a custom error with status code
type CustomError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new custom error
func NewCustomError(message string, statusCode int, err error) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// NewBadRequestError creates a bad request error (400)
func NewBadRequestError(message string) *CustomError {
	return NewCustomError(message, 400, nil)
}

// NewUnauthorizedError creates an unauthorized error (401)
func NewUnauthorizedError(message string) *CustomError {
	return NewCustomError(message, 401, nil)
}

// NewForbiddenError creates a forbidden error (403)
func NewForbiddenError(message string) *CustomError {
	return NewCustomError(message, 403, nil)
}

// NewNotFoundError creates a not found error (404)
func NewNotFoundError(message string) *CustomError {
	return NewCustomError(message, 404, nil)
}

// NewConflictError creates a conflict error (409)
func NewConflictError(message string) *CustomError {
	return NewCustomError(message, 409, nil)
}

// NewInternalError creates an internal server error (500)
func NewInternalError(message string, err error) *CustomError {
	return NewCustomError(message, 500, err)
}

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	var customErr *CustomError
	if errors.As(err, &customErr) {
		return customErr.StatusCode == 404
	}
	return false
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve *ValidationErrors) Error() string {
	return fmt.Sprintf("validation error with %d fields", len(ve.Errors))
}

// NewValidationErrors creates validation errors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: []ValidationError{},
	}
}

// AddError adds a validation error
func (ve *ValidationErrors) AddError(field, message string) *ValidationErrors {
	ve.Errors = append(ve.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
	return ve
}

// HasErrors checks if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}
