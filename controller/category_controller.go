package controller

import (
	"net/http"
	dto "wms/dto/response"
	"wms/models"
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

// UpdateCategory endpoint.
func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var payload models.UpdateCategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	category, err := ctrl.service.UpdateCategory(id, payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, category, "Category berhasil diupdate", nil, http.StatusOK)
}

// DeleteCategory endpoint.
func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.service.DeleteCategory(id)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, nil, "Category berhasil dihapus", nil, http.StatusOK)
}

// CreateCategory endpoint.
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	var payload models.CreateCategoryPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	category, err := ctrl.service.CreateCategory(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, category, "Data berhasil ditambahkan", nil, http.StatusCreated)
}

// ListCategories endpoint.
func (ctrl *CategoryController) ListCategories(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	categories, total, err := ctrl.service.ListCategoriesPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}

	resp := make([]dto.ListCategory, 0, len(categories))

	for _, c := range categories {
		resp = append(resp, dto.ListCategory{
			ID:        c.ID.String(),
			Name:      c.Name,
			Slug:      c.Slug,
			Discount:  c.Discount,
			MinPrice:  (*float64)(c.MinPrice),
			MaxPrice:  (*float64)(c.MaxPrice),
			Status:    c.Status,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			DeletedAt: c.DeletedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	utils.SendListSuccess(c, categories, pg.Page, pg.Limit, total, "", http.StatusOK)
}

func (ctrl *CategoryController) ListCategoriesSelect(c *gin.Context) {
	search := c.Query("search")

	categories, err := ctrl.service.ListCategories(search)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}

	resp := make([]dto.ListCategorySelect, 0, len(categories))
	for _, c := range categories {
		resp = append(resp, dto.ListCategorySelect{
			ID:       c.ID.String(),
			Name:     c.Name,
			Slug:     c.Slug,
			Discount: c.Discount,
			MinPrice: (*float64)(c.MinPrice),
			MaxPrice: (*float64)(c.MaxPrice),
			Status:   c.Status,
		})
	}

	utils.SendItemSuccess(c, resp, "List Category berhasil di ambil", http.StatusOK)
}

// GetCategoryByID endpoint.
func (ctrl *CategoryController) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := ctrl.service.GetCategoryByID(id)
	if err != nil {
		utils.SendError(c, 404, "Category tidak ditemukan")
		return
	}
	utils.SendItemSuccess(c, category, "", http.StatusOK)
}
