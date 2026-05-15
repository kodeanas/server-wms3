package controller

import (
	"net/http"
	dto "wms/dto/response"
	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

// BuyerController defines handlers for buyer resources.
type BuyerController struct {
	service services.BuyerService
}

// NewBuyerController constructor.
func NewBuyerController(service services.BuyerService) *BuyerController {
	return &BuyerController{service: service}
}

// CreateBuyer endpoint.
func (ctrl *BuyerController) CreateBuyer(c *gin.Context) {
	var payload models.CreateBuyerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	buyer, err := ctrl.service.CreateBuyer(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, buyer, "Buyer berhasil ditambahkan", nil, http.StatusCreated)
}

// GetBuyerByID endpoint.
func (ctrl *BuyerController) GetBuyerByID(c *gin.Context) {
	id := c.Param("id")
	detail, err := ctrl.service.GetBuyerDetail(id)
	if err != nil {
		utils.SendError(c, 404, err.Error())
		return
	}
	b := detail.Buyer
	className := ""
	if detail.Class != nil {
		className = detail.Class.Name
	}
	resp := dto.BuyerResponse{
		ID:        b.ID.String(),
		Name:      b.Name,
		Email:     b.Email,
		Phone:     b.Phone,
		ClassName: className,
		Address:   b.Address,
		CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		DeletedAt: b.DeletedAt,
	}
	utils.SendItemSuccess(c, resp, "", http.StatusOK)
}

// ListBuyers endpoint.
func (ctrl *BuyerController) ListBuyers(c *gin.Context) {
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	buyers, total, err := ctrl.service.ListBuyersPaginated(pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	resp := make([]dto.BuyerResponse, 0, len(buyers))
	for _, b := range buyers {
		ClassName := ""
		if b.Class != nil {
			ClassName = b.Class.Name
		}
		resp = append(resp, dto.BuyerResponse{
			ID:        b.ID.String(),
			Name:      b.Name,
			Email:     b.Email,
			Phone:     b.Phone,
			ClassName: ClassName,
			Address:   b.Address,
			CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			DeletedAt: b.DeletedAt,
		})
	}
	utils.SendListSuccess(c, resp, pg.Page, pg.Limit, total, "", http.StatusOK)
}

// UpdateBuyer endpoint.
func (ctrl *BuyerController) UpdateBuyer(c *gin.Context) {
	id := c.Param("id")
	var payload models.UpdateBuyerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}
	buyer, err := ctrl.service.UpdateBuyer(id, payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, buyer, "Buyer berhasil diupdate", nil, http.StatusOK)
}

// DeleteBuyer endpoint.
func (ctrl *BuyerController) DeleteBuyer(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteBuyer(id); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, nil, "Buyer berhasil dihapus", nil, http.StatusOK)
}

// ListBuyersByClass endpoint: list semua buyer berdasarkan class id.
func (ctrl *BuyerController) ListBuyersByClass(c *gin.Context) {
	classID := c.Param("id")
	pg := utils.ParsePagination(c, 10)
	search := c.Query("search")

	buyers, _, total, err := ctrl.service.ListBuyersByClass(classID, pg.Page, pg.Limit, search)
	if err != nil {
		utils.SendError(c, 404, err.Error())
		return
	}
	resp := make([]dto.BuyerClassResponse, 0, len(buyers))
	for _, b := range buyers {
		className := ""
		if b.Class != nil {
			className = b.Class.Name
		}
		resp = append(resp, dto.BuyerClassResponse{
			ID:        b.ID.String(),
			Name:      b.Name,
			Email:     b.Email,
			Phone:     b.Phone,
			ClassName: className,
			Address:   b.Address,
		})
	}

	totalPages := 0
	if pg.Limit > 0 {
		totalPages = int((total + int64(pg.Limit) - 1) / int64(pg.Limit))
	}
	meta := map[string]interface{}{
		"pagination": map[string]interface{}{
			"page":        pg.Page,
			"limit":       pg.Limit,
			"total_items": total,
			"total_pages": totalPages,
		},
	}
	utils.SendSuccess(c, resp, "Data list berhasil diambil", meta, http.StatusOK)
}
