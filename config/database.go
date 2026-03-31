package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"wms/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection with connection pooling
func InitDB() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return err
	}

	// Get underlying SQL database to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
		return err
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("Database connected successfully with connection pooling configured")
	return nil
}

// MigrateDB runs all migrations
func MigrateDB() error {
	return DB.AutoMigrate(
		// User & Auth related
		&models.User{},
		&models.Tax{},

		// Category & Classification
		&models.Category{},
		&models.Sticker{},
		&models.Class{},
		&models.Buyer{},
		&models.UserClassLog{},

		// Products
		&models.ProductDocument{},
		&models.ProductMaster{},
		&models.ProductLog{},

		// Warehouse
		&models.Rack{},
		&models.Store{},
		&models.StoreCrew{},

		// Transfer
		&models.StoreTransfer{},
		&models.StoreTransferBag{},

		// Cargo & Bag
		&models.Cargo{},
		&models.Bag{},

		// Orders
		&models.Order{},
		&models.OrderItem{},
		&models.OrderCargo{},

		// Slow Moving
		&models.SlowMoving{},
		&models.SlowMovingItem{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// GetDatabaseURL returns the PostgreSQL connection URL for migrations
func GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}
