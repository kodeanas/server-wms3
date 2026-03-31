package routes

import (
	"wms/controller"
	"wms/repositories"
	"wms/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r *gin.Engine) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	productRepo := repositories.NewProductRepository()
	orderRepo := repositories.NewOrderRepository()

	// Initialize services
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo)

	// Initialize controllers
	userCtrl := controller.NewUserController(userService)
	productCtrl := controller.NewProductController(productService)
	orderCtrl := controller.NewOrderController(orderService)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "WMS API is running",
		})
	})

	// User routes
	userGroup := r.Group("/api/v1/users")
	{
		userGroup.POST("", userCtrl.CreateUser)
		userGroup.GET("", userCtrl.ListUsers)
		userGroup.GET("/:id", userCtrl.GetUser)
		userGroup.PUT("/:id", userCtrl.UpdateUser)
		userGroup.DELETE("/:id", userCtrl.DeleteUser)
	}

	// Product routes
	productGroup := r.Group("/api/v1/products")
	{
		productGroup.POST("", productCtrl.CreateProduct)
		productGroup.GET("", productCtrl.ListProducts)
		productGroup.GET("/:id", productCtrl.GetProduct)
		productGroup.GET("/barcode/:barcode", productCtrl.GetProductByBarcode)
		productGroup.GET("/category/:categoryID", productCtrl.GetProductsByCategory)
		productGroup.PUT("/:id", productCtrl.UpdateProduct)
		productGroup.DELETE("/:id", productCtrl.DeleteProduct)
	}

	// Order routes
	orderGroup := r.Group("/api/v1/orders")
	{
		orderGroup.POST("", orderCtrl.CreateOrder)
		orderGroup.GET("", orderCtrl.ListOrders)
		orderGroup.GET("/:id", orderCtrl.GetOrder)
		orderGroup.GET("/code/:code", orderCtrl.GetOrderByCode)
		orderGroup.GET("/status/:status", orderCtrl.GetOrdersByStatus)
		orderGroup.PUT("/:id", orderCtrl.UpdateOrder)
		orderGroup.DELETE("/:id", orderCtrl.DeleteOrder)
	}
}
