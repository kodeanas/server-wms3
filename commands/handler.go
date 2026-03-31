package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"wms/config"
	"wms/db"
	"wms/routes"

	"github.com/gin-gonic/gin"
)

// Execute processes CLI commands
func Execute() {
	if len(os.Args) < 2 {
		startServer()
		return
	}

	command := os.Args[1]

	// Load environment first
	config.LoadEnv()

	// Get database URL
	databaseURL := config.GetDatabaseURL()

	// Initialize migration
	m := db.NewMigration(databaseURL)

	switch strings.ToLower(command) {
	case "migrate":
		handleMigrate(m)
	case "migrate:rollback":
		handleRollback(m)
	case "migrate:status":
		handleStatus(m)
	case "migrate:fresh":
		handleFresh(m)
	case "migrate:seed":
		handleSeed(m)
	case "serve":
		startServer()
	default:
		printHelp()
	}
}

// handleMigrate runs all pending migrations
func handleMigrate(m *db.Migration) {
	fmt.Println("🔄 Running migrations...")
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

// handleRollback rolls back the last migration
func handleRollback(m *db.Migration) {
	fmt.Println("⏮️  Rolling back last migration...")
	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
}

// handleStatus shows current migration status
func handleStatus(m *db.Migration) {
	if err := m.Status(); err != nil {
		log.Fatal(err)
	}
}

// handleFresh drops all tables and re-runs migrations
func handleFresh(m *db.Migration) {
	fmt.Println("🔄 Refreshing database (this will drop all tables)...")
	input := ""
	fmt.Print("Are you sure? (yes/no): ")
	fmt.Scanln(&input)

	if strings.ToLower(input) != "yes" {
		fmt.Println("❌ Cancelled")
		return
	}

	if err := m.Fresh(); err != nil {
		log.Fatal(err)
	}
}

// handleSeed seeds the database
func handleSeed(m *db.Migration) {
	if err := m.Seed(); err != nil {
		log.Fatal(err)
	}
}

// startServer starts the HTTP server
func startServer() {
	// Run the main application
	config.LoadEnv()
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := config.MigrateDB(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("🚀 Starting WMS Server...")

	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.Use(config.CORSMiddleware())
	r.Use(config.RequestIDMiddleware())

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("API_HOST")
	if host == "" {
		host = "localhost"
	}

	serverAddr := host + ":" + port
	fmt.Printf("✅ Server running on http://%s\n", serverAddr)
	r.Run(":" + port)
}

// printHelp prints available commands
func printHelp() {
	help := `
╔════════════════════════════════════════════════════════════╗
║              WMS CLI - Available Commands                   ║
╚════════════════════════════════════════════════════════════╝

📋 MIGRATION COMMANDS:
  go run main.go migrate              Run all pending migrations
  go run main.go migrate:rollback     Rollback the last migration
  go run main.go migrate:status       Show current migration status
  go run main.go migrate:fresh        Drop all tables & re-run migrations (⚠️  be careful!)
  go run main.go migrate:seed         Seed the database with sample data

🚀 SERVER COMMANDS:
  go run main.go serve                Start the HTTP server (default)
  go run main.go                      Start the HTTP server (default)

📚 EXAMPLES:
  • First time setup:
    go run main.go migrate
  
  • Check current state:
    go run main.go migrate:status
  
  • Rollback to previous state:
    go run main.go migrate:rollback
  
  • Start development server:
    go run main.go

╚════════════════════════════════════════════════════════════╝
`
	fmt.Println(help)
}
