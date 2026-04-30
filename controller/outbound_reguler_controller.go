package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type OutboundRegulerController struct {
	service services.OutboundRegulerService
}

func NewOutboundRegulerController(service services.OutboundRegulerService) *OutboundRegulerController {
	return &OutboundRegulerController{service: service}
}

// GET /buyers
func (ctrl *OutboundRegulerController) GetBuyers(c *gin.Context) {
	res := ctrl.service.GetBuyers()
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// GET /buyers/:id/class-info
func (ctrl *OutboundRegulerController) GetBuyerClassInfo(c *gin.Context) {
	id := c.Param("id")
	res := ctrl.service.GetBuyerClassInfo(id)
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// POST /outbound-reguler/scan
func (ctrl *OutboundRegulerController) ScanProduct(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	res := ctrl.service.ScanProduct(req)
	if m, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := m["error"]; exists {
			utils.SendError(c, http.StatusBadRequest, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// POST /outbound-reguler/product
func (ctrl *OutboundRegulerController) AddProduct(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	res := ctrl.service.AddProduct(req)
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// DELETE /outbound-reguler/product/:id
func (ctrl *OutboundRegulerController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	res := ctrl.service.DeleteProduct(id)
	if m, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := m["error"]; exists {
			utils.SendError(c, 404, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// POST /outbound-reguler/discount
func (ctrl *OutboundRegulerController) AddDiscount(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	res := ctrl.service.AddDiscount(req)
	if resMap, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := resMap["error"]; exists {
			utils.SendError(c, http.StatusBadRequest, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// PATCH /outbound-reguler/tax
func (ctrl *OutboundRegulerController) UpdateTax(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	res := ctrl.service.UpdateTax(req)
	if resMap, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := resMap["error"]; exists {
			utils.SendError(c, http.StatusBadRequest, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// PATCH /outbound-reguler/box
func (ctrl *OutboundRegulerController) UpdateBox(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	res := ctrl.service.UpdateBox(req)
	if resMap, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := resMap["error"]; exists {
			utils.SendError(c, http.StatusBadRequest, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// POST /outbound-reguler/complete
func (ctrl *OutboundRegulerController) CompleteOrder(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, []utils.ErrorItem{{Field: "", Message: err.Error()}})
		return
	}
	res := ctrl.service.CompleteOrder(req)
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// GET /outbound-reguler/:order_id
func (ctrl *OutboundRegulerController) GetOrderDetail(c *gin.Context) {
	orderID := c.Param("order_id")
	res := ctrl.service.GetOrderDetail(orderID)
	if resMap, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := resMap["error"]; exists {
			utils.SendError(c, http.StatusBadRequest, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "", nil, http.StatusOK)
}

// GET /outbound-reguler/orders
func (ctrl *OutboundRegulerController) ListOrders(c *gin.Context) {
	res := ctrl.service.ListOrders()
	utils.SendSuccess(c, res, "List orders", nil, http.StatusOK)
}

// DELETE /outbound-reguler/discount/order/:order_id
func (ctrl *OutboundRegulerController) DeleteAllDiscountsByOrderID(c *gin.Context) {
	orderID := c.Param("order_id")
	res := ctrl.service.DeleteAllDiscountsByOrderID(orderID)
	if m, ok := res.(map[string]interface{}); ok {
		if errMsg, exists := m["error"]; exists {
			utils.SendError(c, 404, errMsg.(string))
			return
		}
	}
	utils.SendSuccess(c, res, "Voucher/discount berhasil dihapus", nil, http.StatusOK)
}
