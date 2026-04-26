// Tambahkan import jika belum
// import "wms/controller"
// Daftarkan endpoint baru untuk list dokumen SKU
// Misal, jika sudah ada skuController:
// router.GET("/inbound/sku-documents", skuController.ListSKUProductDocuments)
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
	productPendingRepo := repositories.NewProductPendingRepository(config.DB)
	productRepairRepo := repositories.NewProductRepairRepository(config.DB)
	rackDisplayRepo := repositories.NewRackDisplayRepository(config.DB)
	rackStagingRepo := repositories.NewRackStagingRepository(config.DB)
	bagRepo := repositories.NewBagRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)
	taxRepo := repositories.NewTaxRepository(config.DB)
	// TODO: Tambahkan repository product_pending dan product_repair jika sudah ada
	// productPendingRepo := repositories.NewProductPendingRepository(config.DB)
	// productRepairRepo := repositories.NewProductRepairRepository(config.DB)

	// Services
	categoryService := services.NewCategoryService(categoryRepo)
	stickerService := services.NewStickerService(stickerRepo)
	buyerService := services.NewBuyerService(buyerRepo, classRepo)
	classService := services.NewClassService(classRepo)
	productMasterService := services.NewProductMasterService(productMasterRepo)
	productDocumentService := services.NewProductDocumentService(productDocumentRepo)
	productMasterSummaryService := services.NewProductMasterSummaryService(productMasterRepo)
	// Inbound SKU Service
	inboundSKUService := services.NewInboundSKUService(productDocumentRepo, productPendingRepo, productRepairRepo, productMasterRepo)
	rackDisplayService := services.NewRackDisplayService(rackDisplayRepo)
	rackStagingService := services.NewRackStagingService(rackStagingRepo, rackDisplayRepo)
	rackStagingStickerService := services.NewRackStagingStickerService(bagRepo, productMasterRepo)
	wholesaleBagService := services.NewWholesaleBagService(bagRepo, productMasterRepo)
	userService := services.NewUserService(userRepo)
	taxService := services.NewTaxService(taxRepo)

	// Controllers
	categoryController := controller.NewCategoryController(categoryService)
	stickerController := controller.NewStickerController(stickerService)
	buyerController := controller.NewBuyerController(buyerService)
	classController := controller.NewClassController(classService)
	productMasterController := controller.NewProductMasterController(productMasterService)
	productDocumentController := controller.NewProductDocumentController(productDocumentService)
	productMasterSummaryController := controller.NewProductMasterSummaryController(productMasterSummaryService)
	inboundSKUController := controller.NewInboundSKUController(inboundSKUService)
	rackDisplayController := controller.NewRackDisplayController(rackDisplayService)
	rackStagingController := controller.NewRackStagingController(rackStagingService)
	rackStagingStickerController := controller.NewRackStagingStickerController(rackStagingStickerService)
	wholesaleBagController := controller.NewWholesaleBagController(wholesaleBagService)
	userController := controller.NewUserController(userService)
	taxController := controller.NewTaxController(taxService)

	// Public API
	api := r.Group("/api")
	{
		// Taxes
		api.POST("/taxes", taxController.Create)
		api.GET("/taxes", taxController.List)
		api.GET("/taxes/:id", taxController.GetByID)
		api.PUT("/taxes/:id", taxController.Update)
		api.DELETE("/taxes/:id", taxController.Delete)
		api.GET("/taxes-active", taxController.GetActive)

		// Users
		api.POST("/users", userController.CreateUser)
		api.GET("/users", userController.ListUsers)
		api.GET("/users/:id", userController.GetUserByID)
		api.POST("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)
		api.POST("/users/:id/password", userController.UpdatePassword)

		// Rack Displays
		api.POST("/rack-displays", rackDisplayController.Create)
		api.GET("/rack-displays", rackDisplayController.GetAll)
		api.GET("/rack-displays/:id/detail", rackDisplayController.GetByID)
		api.PUT("/rack-displays/:id", rackDisplayController.Update)
		api.DELETE("/rack-displays/:id", rackDisplayController.Delete)

		// Rack Stagings
		api.GET("/rack-stagings", rackStagingController.ListAll)
		api.GET("/rack-stagings/:rackStagingID", rackStagingController.GetDetail)
		api.POST("/rack-stagings", rackStagingController.Create)
		api.GET("/rack-stagings/:rackStagingID/products", productMasterController.ListByRackStagingID)
		api.POST("/rack-stagings/:rackStagingID/scanner/scan-barcode", productMasterController.ScanBarcodeWarehouse)
		api.POST("/rack-stagings/:rackStagingID/finish", rackStagingController.Finish)

		// Rack Staging Sticker (Bag)
		api.POST("/rack-stagings-sticker", rackStagingStickerController.Create)
		api.GET("/rack-stagings-sticker", rackStagingStickerController.List)
		api.GET("/rack-stagings-sticker/:id", rackStagingStickerController.GetDetail)
		api.GET("/rack-stagings-sticker/:id/products", rackStagingStickerController.ListByBagID)
		api.POST("/rack-stagings-sticker/:id/scanner/scan-barcode", rackStagingStickerController.ScanBarcodeWarehouse)

		// Wholesale Bag
		api.POST("/wholesale-bags", wholesaleBagController.Create)
		api.GET("/wholesale-bags", wholesaleBagController.List)
		api.GET("/wholesale-bags/:bagID/detail", wholesaleBagController.GetDetail)
		api.GET("/wholesale-bags/:bagID/products", wholesaleBagController.ListByBagID)
		api.POST("/wholesale-bags/:bagID/scanner/scan-barcode", wholesaleBagController.ScanBarcodeWarehouse)

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
			   api.PUT("/classes/:id/up", classController.MoveUp)
			   api.PUT("/classes/:id/down", classController.MoveDown)

		// Inbound Manual
		api.GET("/inbound/list-masters", controller.ListAllProductMastersHandler(config.DB))
		api.GET("/inbound/list-pendings", controller.ListAllProductPendingsHandler(config.DB))

		api.GET("/inbound/manual-pending", controller.ListProductManualHandler(config.DB))
		api.POST("/inbound/manual", controller.InboundManualHandler(config.DB))
		api.POST("/inbound/bulk-upload", controller.InboundBulkUploadHandler(config.DB))
		api.POST("/inbound/bast-upload", controller.InboundBastUploadHandler(config.DB))

		// Inbound BAST
		api.GET("/inbound/bast-summary", controller.InboundBastSummaryHandler(config.DB))

		api.GET("/inbound/bast-scanner/document/:document_id", controller.InboundBastGetDocumentHandler(config.DB))
		api.GET("/inbound/bast-scanner/:document_id/product/:barcode", controller.InboundBastGetPendingProductHandler(config.DB))
		api.POST("/inbound/bast-scanner/:document_id/scan/:barcode", controller.InboundBastScanSingleProductHandler(config.DB))
		api.POST("/inbound/bast-scanner/:document_id/finish", productDocumentController.FinishDocument)

		// Product Master Staging Reguler
		api.GET("/product-masters/staging-reguler", productMasterController.ListStagingReguler)
		api.GET("/product-masters/staging-sticker", productMasterController.ListStagingSticker)
		api.POST("/product-masters/staging/:id", productMasterController.UpdateStaging)
		api.GET("/product-masters/:id", productMasterController.GetDetail)

		// Product Document
		api.GET("/product-documents", productDocumentController.ListDocuments)
		api.GET("/product-documents/bulk", productDocumentController.GetBulkDocuments)
		api.GET("/product-documents/bulk/:id", productDocumentController.GetBulkDocumentDetail)
		api.GET("/product-documents/bast", productDocumentController.GetBastDocuments)
		api.GET("/product-documents/bast/:id/relations", productDocumentController.GetBastRelationsDetail)
		api.GET("/product-documents/bast/:id/overview", productDocumentController.GetBastOverview)
		api.GET("/product-documents/bast/:id/pending-by-type", productDocumentController.GetBastPendingByType)
		api.GET("/product-documents/sku", inboundSKUController.ListSKUProductDocuments)

		// Product Master Summary
		api.GET("/manual/summary", productMasterSummaryController.GetSummary)

		// Inbound SKU
		api.POST("/inbound-sku/upload", inboundSKUController.UploadExcel)
		api.POST("/inbound-sku/crosscheck/:pending_id", inboundSKUController.CrosscheckPending)
		api.POST("/inbound-sku/finish/:document_id", inboundSKUController.FinishInboundSKU)
		api.GET("/inbound-sku/document/:document_id", controller.InboundSKUGetDocumentHandler(config.DB))
	}
}
