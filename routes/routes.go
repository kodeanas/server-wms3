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
	// Category
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryController := controller.NewCategoryController(categoryService)

	//Sticker
	stickerRepo := repositories.NewStickerRepository(config.DB)
	stickerService := services.NewStickerService(stickerRepo)
	stickerController := controller.NewStickerController(stickerService)

	//Class
	classRepo := repositories.NewClassRepository(config.DB)
	classService := services.NewClassService(classRepo)
	classController := controller.NewClassController(classService)

	//Buyer
	buyerRepo := repositories.NewBuyerRepository(config.DB)
	buyerService := services.NewBuyerService(buyerRepo, classRepo)
	buyerController := controller.NewBuyerController(buyerService)

	//User
	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	//Tax
	taxRepo := repositories.NewTaxRepository(config.DB)
	taxService := services.NewTaxService(taxRepo)
	taxController := controller.NewTaxController(taxService)

	//Product Document
	productDocumentRepo := repositories.NewProductDocumentRepository(config.DB)
	productDocumentService := services.NewProductDocumentService(productDocumentRepo)
	productDocumentController := controller.NewProductDocumentController(productDocumentService)

	//Product Master
	productMasterRepo := repositories.NewProductMasterRepository(config.DB)
	productMasterService := services.NewProductMasterService(productMasterRepo)
	productMasterController := controller.NewProductMasterController(productMasterService)

	//Product Master Summary
	productMasterSummaryService := services.NewProductMasterSummaryService(productMasterRepo)
	productMasterSummaryController := controller.NewProductMasterSummaryController(productMasterSummaryService)

	//Product Pending & Product Repair
	productPendingRepo := repositories.NewProductPendingRepository(config.DB)
	productRepairRepo := repositories.NewProductRepairRepository(config.DB)

	//Inbound SKU
	inboundSKUService := services.NewInboundSKUService(productDocumentRepo, productPendingRepo, productRepairRepo, productMasterRepo)
	inboundSKUController := controller.NewInboundSKUController(inboundSKUService)

	//RackDisplay
	rackDisplayRepo := repositories.NewRackDisplayRepository(config.DB)
	rackDisplayService := services.NewRackDisplayService(rackDisplayRepo)
	rackDisplayController := controller.NewRackDisplayController(rackDisplayService)

	//RWholesale Bag
	bagRepo := repositories.NewBagRepository(config.DB)
	BagService := services.NewWholesaleBagService(bagRepo, productMasterRepo)
	BagController := controller.NewWholesaleBagController(BagService)

	//Rack Staging Reguler
	rackStagingRepo := repositories.NewRackStagingRepository(config.DB)
	rackStagingRegulerService := services.NewRackStagingService(rackStagingRepo, rackDisplayRepo)
	rackStagingRegulerController := controller.NewRackStagingController(rackStagingRegulerService)

	//Rack Staging Sticker
	rackStagingStickerService := services.NewRackStagingStickerService(bagRepo, productMasterRepo)
	rackStagingStickerController := controller.NewRackStagingStickerController(rackStagingStickerService)

	//Order ,ProductOrder , DiscountOrder
	orderRepo := repositories.NewOrderRepository(config.DB)
	productOrderRepo := repositories.NewProductOrderRepository(config.DB)
	discountOrderRepo := repositories.NewDiscountOrderRepository(config.DB)

	// Outbound Reguler
	outboundRegulerService := services.NewOutboundRegulerService(buyerRepo, classRepo, orderRepo, productOrderRepo, discountOrderRepo, categoryRepo, productMasterRepo, taxRepo)
	outboundRegulerController := controller.NewOutboundRegulerController(outboundRegulerService)

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
		api.PUT("/classes/:id/up", classController.MoveUp)
		api.PUT("/classes/:id/down", classController.MoveDown)

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
		api.PUT("/users/:id", userController.UpdateUser)
		api.DELETE("/users/:id", userController.DeleteUser)
		api.PUT("/users/:id/password", userController.UpdatePassword)

		// Product Master Summary
		api.GET("/manual/summary", productMasterSummaryController.GetSummary)

		// Product Master Staging Reguler
		api.GET("/product-masters/staging-reguler", productMasterController.ListStagingReguler)
		api.GET("/product-masters/staging-sticker", productMasterController.ListStagingSticker)
		api.POST("/product-masters/staging/:id", productMasterController.UpdateStaging)
		api.GET("/product-masters/:id", productMasterController.GetDetail)

		// Product Document
		api.GET("/product-documents", productDocumentController.ListDocuments)
		api.GET("/product-documents/bulk", productDocumentController.GetBulkDocuments)
		api.GET("/product-documents/bulk/:id", productDocumentController.GetBulkDocumentDetail)

		api.GET("/product-documents/bast/:id/relations", productDocumentController.GetBastRelationsDetail)
		api.GET("/product-documents/bast/:id/overview", productDocumentController.GetBastOverview)
		api.GET("/product-documents/bast/:id/pending-by-type", productDocumentController.GetBastPendingByType)
		api.GET("/product-documents/sku", inboundSKUController.ListSKUProductDocuments)

		// Inbound Manual
		api.GET("/inbound/list-masters", controller.ListAllProductMastersHandler(config.DB))
		api.GET("/inbound/list-pendings", controller.ListAllProductPendingsHandler(config.DB))
		api.GET("/inbound/manual-pending", controller.ListProductManualHandler(config.DB))

		api.POST("/inbound/manual", controller.InboundManualHandler(config.DB))

		//Inbound Bulk
		api.POST("/inbound/bulk-upload", controller.InboundBulkUploadHandler(config.DB))

		// Inbound BAST
		api.POST("/inbound/bast-upload", controller.InboundBastUploadHandler(config.DB))
		api.GET("/inbound/bast-summary", controller.InboundBastSummaryHandler(config.DB))

		api.GET("/product-documents/bast", productDocumentController.GetBastDocuments)
		api.GET("/inbound/bast-scanner/document/:document_id", controller.InboundBastGetDocumentHandler(config.DB))
		api.GET("/inbound/bast-scanner/:document_id/product/:barcode", controller.InboundBastGetPendingProductHandler(config.DB))
		api.POST("/inbound/bast-scanner/:document_id/scan/:barcode", controller.InboundBastScanSingleProductHandler(config.DB))
		api.POST("/inbound/bast-scanner/:document_id/finish", productDocumentController.FinishDocument)

		// Inbound SKU
		api.POST("/inbound-sku/upload", inboundSKUController.UploadExcel)
		api.POST("/inbound-sku/crosscheck/:pending_id", inboundSKUController.CrosscheckPending)
		api.POST("/inbound-sku/finish/:document_id", inboundSKUController.FinishInboundSKU)
		api.GET("/inbound-sku/document/:document_id", controller.InboundSKUGetDocumentHandler(config.DB))

		// Rack Displays
		api.POST("/rack-displays", rackDisplayController.Create)
		api.GET("/rack-displays", rackDisplayController.GetAll)
		api.GET("/rack-displays/:id/detail", rackDisplayController.GetDetail)
		api.PUT("/rack-displays/:id", rackDisplayController.Update)
		api.DELETE("/rack-displays/:id", rackDisplayController.Delete)

		// Rack Stagings Reguler
		api.GET("/rack-stagings", rackStagingRegulerController.ListAll)
		api.GET("/rack-stagings/:rackStagingID", rackStagingRegulerController.GetDetail)
		api.POST("/rack-stagings", rackStagingRegulerController.Create)
		api.GET("/rack-stagings/:rackStagingID/products", productMasterController.ListByRackStagingID)
		api.POST("/rack-stagings/:rackStagingID/scanner/scan-barcode", productMasterController.ScanBarcodeWarehouse)
		api.POST("/rack-stagings/:rackStagingID/finish", rackStagingRegulerController.Finish)

		// Rack Staging Sticker (Bag)
		api.POST("/rack-stagings-sticker", rackStagingStickerController.Create)
		api.GET("/rack-stagings-sticker", rackStagingStickerController.List)
		api.GET("/rack-stagings-sticker/:id", rackStagingStickerController.GetDetail)
		api.GET("/rack-stagings-sticker/:id/products", rackStagingStickerController.ListByBagID)
		api.POST("/rack-stagings-sticker/:id/scanner/scan-barcode", rackStagingStickerController.ScanBarcodeWarehouse)

		// Wholesale Bag
		api.POST("/wholesale-bags", BagController.Create)
		api.GET("/wholesale-bags", BagController.List)
		api.GET("/wholesale-bags/:bagID/detail", BagController.GetDetail)
		api.GET("/wholesale-bags/:bagID/products", BagController.ListByBagID)
		api.POST("/wholesale-bags/:bagID/scanner/scan-barcode", BagController.ScanBarcodeWarehouse)

		// Outbound Reguler
		api.GET("/outbound-reguler/buyers", outboundRegulerController.GetBuyers)
		api.GET("/outbound-reguler/buyers/:id/class-info", outboundRegulerController.GetBuyerClassInfo)
		api.POST("/outbound-reguler/scan", outboundRegulerController.ScanProduct)
		api.DELETE("/outbound-reguler/product/:id", outboundRegulerController.DeleteProduct)
		api.POST("/outbound-reguler/discount", outboundRegulerController.AddDiscount)
		api.DELETE("/outbound-reguler/discount/order/:order_id", outboundRegulerController.DeleteAllDiscountsByOrderID)
		api.PATCH("/outbound-reguler/tax", outboundRegulerController.UpdateTax)
		api.PATCH("/outbound-reguler/box", outboundRegulerController.UpdateBox)
		api.POST("/outbound-reguler/complete", outboundRegulerController.CompleteOrder)
		api.GET("/outbound-reguler/orders", outboundRegulerController.ListOrders)
		api.GET("/outbound-reguler/:order_id", outboundRegulerController.GetOrderDetail)
	}
}
