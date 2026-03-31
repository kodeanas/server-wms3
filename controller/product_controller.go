package controller

import (
	"net/http"
	"strconv"

	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

// NewProductController creates a new product controller
func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{service: service}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product in the system
// @Accept json
// @Produce json
// @Param product body models.ProductMaster true "Product data"
// @Success 201 {object} utils.Response
// @Router /products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product models.ProductMaster
	if err := ctx.ShouldBindJSON(&product); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.service.CreateProduct(ctx.Request.Context(), &product); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.SendSuccess(ctx, product, "Product created successfully", http.StatusCreated)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Retrieve a product by their ID
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response
// @Router /products/{id} [get]
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.service.GetProduct(ctx.Request.Context(), id)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Product not found")
		return
	}

	utils.SendSuccess(ctx, product, "Product retrieved successfully")
}

// GetProductByBarcode godoc
// @Summary Get product by barcode
// @Description Retrieve a product by barcode
// @Param barcode path string true "Product barcode"
// @Success 200 {object} utils.Response
// @Router /products/barcode/{barcode} [get]
func (c *ProductController) GetProductByBarcode(ctx *gin.Context) {
	barcode := ctx.Param("barcode")

	product, err := c.service.GetProductByBarcode(ctx.Request.Context(), barcode)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Product not found")
		return
	}

	utils.SendSuccess(ctx, product, "Product retrieved successfully")
}

// ListProducts godoc
// @Summary List all products
// @Description Retrieve paginated list of products
// @Param limit query int false "Limit results" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} utils.Response
// @Router /products [get]
func (c *ProductController) ListProducts(ctx *gin.Context) {
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

	products, total, err := c.service.ListProducts(ctx.Request.Context(), limit, offset)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"products": products,
			"total":    total,
		},
		"message": "Products retrieved successfully",
	})
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Retrieve products filtered by category
// @Param categoryID path string true "Category ID"
// @Success 200 {object} utils.Response
// @Router /products/category/{categoryID} [get]
func (c *ProductController) GetProductsByCategory(ctx *gin.Context) {
	categoryID := ctx.Param("categoryID")
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

	products, total, err := c.service.GetProductsByCategory(ctx.Request.Context(), categoryID, limit, offset)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"products": products,
			"total":    total,
		},
		"message": "Products retrieved successfully",
	})
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product
// @Param id path string true "Product ID"
// @Param product body models.ProductMaster true "Product data"
// @Success 200 {object} utils.Response
// @Router /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.ProductMaster

	if err := ctx.ShouldBindJSON(&product); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.service.UpdateProduct(ctx.Request.Context(), id, &product); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to update product")
		return
	}

	utils.SendSuccess(ctx, product, "Product updated successfully")
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product from the system
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response
// @Router /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.service.DeleteProduct(ctx.Request.Context(), id); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	utils.SendSuccess(ctx, nil, "Product deleted successfully")
}
