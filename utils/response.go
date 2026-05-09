package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse defines standard API response.
type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorItem for validation errors.
type ErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// PaginationInfo represents pagination metadata payload.
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// MetaPagination wraps PaginationInfo into the meta envelope.
type MetaPagination struct {
	Pagination PaginationInfo `json:"pagination"`
}

// BuildMetaPagination builds the standardized meta payload for paginated GET LIST responses.
func BuildMetaPagination(page, limit int, totalItems int64) MetaPagination {
	totalPages := 0
	if limit > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}
	return MetaPagination{
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}

// SendSuccess sends a standard success response with optional meta.
// Backward compatible: meta dapat di-pass nil untuk GET satuan.
func SendSuccess(c *gin.Context, data interface{}, message string, meta interface{}, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	if message == "" {
		message = "success"
	}

	// Jika data nil, jadikan map kosong agar di JSON menjadi {}
	respData := data
	if data == nil {
		respData = map[string]interface{}{}
	}
	c.JSON(code, APIResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    respData,
		Meta:    meta,
	})
}

// SendSuccessWithMetaNull is for backward compatibility, always sends meta as null.
func SendSuccessWithMetaNull(c *gin.Context, data interface{}, message string, statusCode ...int) {
	SendSuccess(c, data, message, nil, statusCode...)
}

// SendItemSuccess untuk response GET satuan (tanpa meta).
// Default message: "Data berhasil diambil".
func SendItemSuccess(c *gin.Context, data interface{}, message string, statusCode ...int) {
	if message == "" {
		message = "Data berhasil diambil"
	}
	SendSuccess(c, data, message, nil, statusCode...)
}

// SendListSuccess untuk response GET LIST dengan pagination meta.
// Default message: "Data list berhasil diambil".
func SendListSuccess(c *gin.Context, data interface{}, page, limit int, totalItems int64, message string, statusCode ...int) {
	if message == "" {
		message = "Data list berhasil diambil"
	}
	// Jika data nil, kembalikan array kosong agar JSON-nya [] bukan null.
	respData := data
	if data == nil {
		respData = []interface{}{}
	}
	meta := BuildMetaPagination(page, limit, totalItems)
	SendSuccess(c, respData, message, meta, statusCode...)
}

// SendPaginatedSuccess sends paginated success response.
// Deprecated: gunakan SendListSuccess. Disimpan untuk backward compatibility.
func SendPaginatedSuccess(c *gin.Context, data interface{}, page, limit, totalItems, totalPages int64, message string) {
	if message == "" {
		message = "Data list berhasil diambil"
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Code:    http.StatusOK,
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
		Status:  "error",
		Code:    code,
		Message: message,
		Data:    nil,
		Meta:    nil,
	})
}

// SendValidationError sends validation error response.
func SendValidationError(c *gin.Context, errors []ErrorItem) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Status:  "error",
		Code:    http.StatusBadRequest,
		Message: "Validasi gagal",
		Data:    nil,
		Meta: map[string]interface{}{
			"errors": errors,
		},
	})
}
