package routes

import (
	"wms/config"
	"wms/controller"
	"wms/repositories"
	"wms/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes menginisialisasi semua endpoint pada router yang diberikan.
// ini di-set di sini karena akan semakin banyak resource.
func SetupRoutes(r *gin.Engine) {
	// Repositories
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	stickerRepo := repositories.NewStickerRepository(config.DB)
	buyerRepo := repositories.NewBuyerRepository(config.DB)
	classRepo := repositories.NewClassRepository(config.DB)
	productMasterRepo := repositories.NewProductMasterRepository(config.DB)
	productDocumentRepo := repositories.NewProductDocumentRepository(config.DB)

	// Services
	categoryService := services.NewCategoryService(categoryRepo)
	stickerService := services.NewStickerService(stickerRepo)
	buyerService := services.NewBuyerService(buyerRepo, classRepo)
	classService := services.NewClassService(classRepo)
	productMasterService := services.NewProductMasterService(productMasterRepo)
	productDocumentService := services.NewProductDocumentService(productDocumentRepo)
	productMasterSummaryService := services.NewProductMasterSummaryService(productMasterRepo)

	// Controllers
	categoryController := controller.NewCategoryController(categoryService)
	stickerController := controller.NewStickerController(stickerService)
	buyerController := controller.NewBuyerController(buyerService)
	classController := controller.NewClassController(classService)
	productMasterController := controller.NewProductMasterController(productMasterService)
	productDocumentController := controller.NewProductDocumentController(productDocumentService)
	productMasterSummaryController := controller.NewProductMasterSummaryController(productMasterSummaryService)

	// Public API
	api := r.Group("/api")
	{
		// Categories
		api.POST("/categories", categoryController.CreateCategory)
		api.GET("/categories", categoryController.ListCategories)
		api.GET("/categories/:id", categoryController.GetCategoryByID)
		api.PUT("/categories/:id", categoryController.UpdateCategory)
		api.DELETE("/categories/:id", categoryController.DeleteCategory)

		// Stickers
		api.POST("/stickers", stickerController.CreateSticker)
		api.GET("/stickers", stickerController.ListStickers)
		api.GET("/stickers/:id", stickerController.GetStickerByID)
		api.PUT("/stickers/:id", stickerController.UpdateSticker)
		api.DELETE("/stickers/:id", stickerController.DeleteSticker)

		// Buyers
		api.POST("/buyers", buyerController.CreateBuyer)
		api.GET("/buyers", buyerController.ListBuyers)
		api.GET("/buyers/:id", buyerController.GetBuyerByID)
		api.PUT("/buyers/:id", buyerController.UpdateBuyer)
		api.DELETE("/buyers/:id", buyerController.DeleteBuyer)

		// Classes
		api.POST("/classes", classController.CreateClass)
		api.GET("/classes", classController.ListClasses)
		api.GET("/classes/:id", classController.GetClassByID)
		api.PUT("/classes/:id", classController.UpdateClass)
		api.DELETE("/classes/:id", classController.DeleteClass)

		// Inbound Manual
		api.POST("/scanin/manual", controller.InboundManualHandler(config.DB))
		// api.GET("/scanin/manual", controller.ListAllProductMastersHandler(config.DB))
		api.GET("/scanin/manual", controller.ListProductManualHandler(config.DB))

		// Inbound Bulk (single API)
		api.POST("/inbound/bulk-upload", controller.InboundBulkUploadHandler(config.DB))

		// Product Master Staging Reguler
		api.GET("/product-masters/staging-reguler", productMasterController.ListStagingReguler)
		api.GET("/product-masters/staging-sticker", productMasterController.ListStagingSticker)

		// Product Document
		api.GET("/product-documents", productDocumentController.ListDocuments)

		// Product Master Summary
		api.GET("/manual/summary", productMasterSummaryController.GetSummary)
	}
}
