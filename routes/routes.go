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

	// Services
	categoryService := services.NewCategoryService(categoryRepo)

	// Controllers
	categoryController := controller.NewCategoryController(categoryService)

	// Public API
	api := r.Group("/api")
	{
		api.POST("/categories", categoryController.CreateCategory)
		api.GET("/categories", categoryController.ListCategories)
	}

	// contoh group lain (auth / protected) disiapkan untuk scale
	// auth := r.Group("/auth")
	// {
	// 	auth.POST("/login", authController.Login)
	// 	auth.POST("/register", authController.Register)
	// }

	// protected := r.Group("/api")
	// protected.Use(authMiddleware)
	// {
	// 	protected.GET("/profile", profileController.GetProfile)
	// }
}
