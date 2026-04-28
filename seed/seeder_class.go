package seed

import (
	"fmt"
	"wms/config"
	"wms/models"
)

// SeedClass seeds the database with default class data.
func SeedClass() {
	classes := []models.Class{
		{
			Name:                "Bronze",
			MinOrder:            0,
			Disc:                5,
			MinTransactionValue: 5000000,
			Week:                8,
			Iteration:           1,
			Status:              "active",
		},
		{
			Name:                "Gold",
			MinOrder:            4,
			Disc:                10,
			MinTransactionValue: 5000000,
			Week:                4,
			Iteration:           2,
			Status:              "active",
		},
		{
			Name:                "Diamond",
			MinOrder:            8,
			Disc:                20,
			MinTransactionValue: 5000000,
			Week:                4,
			Iteration:           3,
			Status:              "active",
		},
	}

	for _, class := range classes {
		if err := config.DB.Create(&class).Error; err != nil {
			fmt.Printf("Failed to insert class %s: %v\n", class.Name, err)
		} else {
			fmt.Printf("Inserted class: %s\n", class.Name)
		}
	}
}
