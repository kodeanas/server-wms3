package seed

import (
	"fmt"
	"wms/config"
	"wms/models"
	"wms/utils"
)

// SeedSticker seeds the database with default sticker data.
func SeedSticker() {
	stickers := []models.Sticker{
		{
			CodeHex:    "FF5733",
			Name:       "Sticker A",
			Slug:       utils.GenerateSlug("Sticker A"),
			Type:       "big",
			FixedPrice: IntPtr(75000),
			MinPrice:   PricePtr(50000),
			MaxPrice:   PricePtr(99999),
		},
		{
			CodeHex:    "FFFF00",
			Name:       "Sticker B",
			Slug:       utils.GenerateSlug("Sticker B"),
			Type:       "tiny",
			FixedPrice: IntPtr(30000),
			MinPrice:   PricePtr(20000),
			MaxPrice:   PricePtr(49999),
		},
		{
			CodeHex:    "00FF00",
			Name:       "Sticker C",
			Slug:       utils.GenerateSlug("Sticker C"),
			Type:       "small",
			FixedPrice: IntPtr(12000),
			MinPrice:   PricePtr(0),
			MaxPrice:   PricePtr(19999),
		},
	}

	for _, sticker := range stickers {
		if err := config.DB.Create(&sticker).Error; err != nil {
			fmt.Printf("Failed to insert sticker %s: %v\n", sticker.Name, err)
		} else {
			fmt.Printf("Inserted sticker: %s\n", sticker.Name)
		}
	}
}

// IntPtr returns a pointer to an int value
func IntPtr(i int) *int {
	return &i
}

// PricePtr returns a pointer to a models.Price value
func PricePtr(i int) *models.Price {
	p := models.Price(float64(i))
	return &p
}
