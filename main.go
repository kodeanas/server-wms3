package main

import (
	"log"
	"os"

	"wms/config"
	"wms/models"
	"wms/routes"
	"wms/seed"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seed.SeedSticker()
		seed.SeedCategory()
		seed.SeedClass()
		return // keluar setelah seed
	}

	// Auto migrate models
	if err := config.DB.AutoMigrate(&models.Category{}, &models.Sticker{}); err != nil {
		log.Fatal("auto migrate failed:", err)
	}

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
