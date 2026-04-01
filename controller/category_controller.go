package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

// CategoryController defines handlers for category resources.
type CategoryController struct {
	service services.CategoryService
}

// NewCategoryController constructor.
func NewCategoryController(service services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// CreateCategory endpoint.
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	var payload services.CreateCategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		// if errs, ok := err.(gin.Errors); ok {
		// 	validationErrors = make([]utils.ErrorItem, 0, len(errs))
		// 	for _, e := range errs {
		// 		validationErrors = append(validationErrors, utils.ErrorItem{Field: e.Field(), Message: e.Error()})
		// 	}
		// }
		utils.SendValidationError(c, validationErrors)
		return
	}

	category, err := ctrl.service.CreateCategory(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, category, "Data berhasil ditambahkan", http.StatusCreated)
}

// ListCategories endpoint.
func (ctrl *CategoryController) ListCategories(c *gin.Context) {
	categories, err := ctrl.service.ListCategories()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, categories, "list categories")
}
