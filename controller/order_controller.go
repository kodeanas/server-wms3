package controller

import (
	"net/http"
	"strconv"

	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	service services.OrderService
}

// NewOrderController creates a new order controller
func NewOrderController(service services.OrderService) *OrderController {
	return &OrderController{service: service}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order in the system
// @Accept json
// @Produce json
// @Param order body models.Order true "Order data"
// @Success 201 {object} utils.Response
// @Router /orders [post]
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var order models.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.service.CreateOrder(ctx.Request.Context(), &order); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to create order")
		return
	}

	utils.SendSuccess(ctx, order, "Order created successfully", http.StatusCreated)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Retrieve an order by their ID
// @Param id path string true "Order ID"
// @Success 200 {object} utils.Response
// @Router /orders/{id} [get]
func (c *OrderController) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.service.GetOrder(ctx.Request.Context(), id)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Order not found")
		return
	}

	utils.SendSuccess(ctx, order, "Order retrieved successfully")
}

// GetOrderByCode godoc
// @Summary Get order by code
// @Description Retrieve an order by order code
// @Param code path string true "Order code"
// @Success 200 {object} utils.Response
// @Router /orders/code/{code} [get]
func (c *OrderController) GetOrderByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	order, err := c.service.GetOrderByCode(ctx.Request.Context(), code)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Order not found")
		return
	}

	utils.SendSuccess(ctx, order, "Order retrieved successfully")
}

// ListOrders godoc
// @Summary List all orders
// @Description Retrieve paginated list of orders
// @Param limit query int false "Limit results" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} utils.Response
// @Router /orders [get]
func (c *OrderController) ListOrders(ctx *gin.Context) {
	limit := 10
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	orders, total, err := c.service.ListOrders(ctx.Request.Context(), limit, offset)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"orders": orders,
			"total":  total,
		},
		"message": "Orders retrieved successfully",
	})
}

// GetOrdersByStatus godoc
// @Summary Get orders by status
// @Description Retrieve orders filtered by status
// @Param status path string true "Order status"
// @Success 200 {object} utils.Response
// @Router /orders/status/{status} [get]
func (c *OrderController) GetOrdersByStatus(ctx *gin.Context) {
	status := ctx.Param("status")
	limit := 10
	offset := 0

	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}

	if o := ctx.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	orders, total, err := c.service.GetOrdersByStatus(ctx.Request.Context(), status, limit, offset)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"orders": orders,
			"total":  total,
		},
		"message": "Orders retrieved successfully",
	})
}

// UpdateOrder godoc
// @Summary Update order
// @Description Update an existing order
// @Param id path string true "Order ID"
// @Param order body models.Order true "Order data"
// @Success 200 {object} utils.Response
// @Router /orders/{id} [put]
func (c *OrderController) UpdateOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.Order

	if err := ctx.ShouldBindJSON(&order); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.service.UpdateOrder(ctx.Request.Context(), id, &order); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to update order")
		return
	}

	utils.SendSuccess(ctx, order, "Order updated successfully")
}

// DeleteOrder godoc
// @Summary Delete order
// @Description Delete an order from the system
// @Param id path string true "Order ID"
// @Success 200 {object} utils.Response
// @Router /orders/{id} [delete]
func (c *OrderController) DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteOrder(ctx.Request.Context(), id); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to delete order")
		return
	}

	utils.SendSuccess(ctx, nil, "Order deleted successfully")
}
