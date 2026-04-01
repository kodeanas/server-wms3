package main

import (
	"log"
	"os"

	"wms/config"
	"wms/controller"
	"wms/models"
	"wms/repositories"
	"wms/routes"
	"wms/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// Auto migrate category table (and related models if desired)
	if err := config.DB.AutoMigrate(&models.Category{}); err != nil {
		log.Fatal("auto migrate failed:", err)
	}

	categoryRepo := repositories.NewCategoryRepository(config.DB)
	categoryService := services.NewCategoryService(categoryRepo)
	_ = controller.NewCategoryController(categoryService) // keep initialization if needed later

	router := gin.Default()

	// optional simple CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
