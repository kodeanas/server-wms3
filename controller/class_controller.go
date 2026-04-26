package controller

import (
	"fmt"
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

// ClassController defines handlers for class resources.
type ClassController struct {
	service services.ClassService
}

// NewClassController constructor.
func NewClassController(service services.ClassService) *ClassController {
	return &ClassController{service: service}
}

// CreateClass endpoint.
func (ctrl *ClassController) CreateClass(c *gin.Context) {
	var payload services.CreateClassPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	class, err := ctrl.service.CreateClass(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, class, "Class berhasil ditambahkan", nil, http.StatusCreated)
}

// GetClassByID endpoint.
func (ctrl *ClassController) GetClassByID(c *gin.Context) {
	id := c.Param("id")
	class, err := ctrl.service.GetClassByID(id)
	if err != nil {
		utils.SendError(c, 404, err.Error())
		return
	}
	resp := map[string]interface{}{
		"id":                    class.ID,
		"name":                  class.Name,
		"min_order":             class.MinOrder,
		"disc":                  class.Disc,
		"min_transaction_value": fmt.Sprintf("%.2f", class.MinTransactionValue),
		"week":                  class.Week,
		"iteration":             class.Iteration,
		"status":                class.Status,
		"created_at":            class.CreatedAt,
		"updated_at":            class.UpdatedAt,
	}
	utils.SendSuccess(c, resp, "", nil, http.StatusOK)
}

// ListClasses endpoint.
func (ctrl *ClassController) ListClasses(c *gin.Context) {
	classes, err := ctrl.service.ListClasses()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	var resp []map[string]interface{}
	for _, class := range classes {
		resp = append(resp, map[string]interface{}{
			"id":                    class.ID,
			"name":                  class.Name,
			"min_order":             class.MinOrder,
			"disc":                  class.Disc,
			"min_transaction_value": fmt.Sprintf("%.2f", class.MinTransactionValue),
			"week":                  class.Week,
			"iteration":             class.Iteration,
			"status":                class.Status,
			"created_at":            class.CreatedAt,
			"updated_at":            class.UpdatedAt,
		})
	}
	utils.SendSuccess(c, resp, "", nil, http.StatusOK)
}

// UpdateClass endpoint.
func (ctrl *ClassController) UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var payload services.UpdateClassPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}
	class, err := ctrl.service.UpdateClass(id, payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, class, "Class berhasil diupdate", nil, http.StatusOK)
}

// DeleteClass endpoint.
func (ctrl *ClassController) DeleteClass(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteClass(id); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, nil, "Class berhasil dihapus", nil, http.StatusOK)
}

// MoveUp endpoint: naikkan urutan class
func (ctrl *ClassController) MoveUp(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.MoveUp(id); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(c, nil, "Class berhasil dinaikkan")
}

// MoveDown endpoint: turunkan urutan class
func (ctrl *ClassController) MoveDown(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.MoveDown(id); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccessWithMetaNull(c, nil, "Class berhasil diturunkan")
}
