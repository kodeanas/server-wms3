package seed

import (
	"fmt"
	"wms/config"
	"wms/models"
	"wms/utils"
)

// SeedCategory seeds the database with default category data.
func SeedCategory() {
	categories := []models.Category{
		{
			Name:     "TOYS HOBBIES (>700)",
			Discount: IntPtr(40),
			MinPrice: PricePtr(100000),
			MaxPrice: PricePtr(10000000),
			Status:   "active",
			Slug:     utils.GenerateSlug("TOYS HOBBIES (>700)"),
		},
		{
			Name:     "TOYS HOBBIES (200-699)",
			Discount: IntPtr(50),
			MinPrice: PricePtr(100000),
			MaxPrice: PricePtr(699999),
			Status:   "active",
			Slug:     utils.GenerateSlug("TOYS HOBBIES (200-699)"),
		},
		{
			Name:     "FMCG",
			Discount: IntPtr(50),
			MinPrice: PricePtr(100000),
			MaxPrice: PricePtr(10000000),
			Status:   "active",
			Slug:     utils.GenerateSlug("FMCG"),
		},
	}

	for _, category := range categories {
		if err := config.DB.Create(&category).Error; err != nil {
			fmt.Printf("Failed to insert category %s: %v\n", category.Name, err)
		} else {
			fmt.Printf("Inserted category: %s\n", category.Name)
		}
	}
}
