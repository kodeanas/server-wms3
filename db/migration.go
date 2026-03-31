package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migration handles database migrations
type Migration struct {
	databaseURL    string
	migrationsPath string
}

// NewMigration creates a new migration instance
func NewMigration(databaseURL string) *Migration {
	// Get absolute path to migrations folder
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path:", err)
	}
	exePath := filepath.Dir(ex)

	// Try relative path first (for development)
	migPath := "db/migrations"
	if _, err := os.Stat(migPath); os.IsNotExist(err) {
		// Try from executable path (for production)
		migPath = filepath.Join(exePath, "db/migrations")
	}

	// Get absolute path
	absPath, err := filepath.Abs(migPath)
	if err != nil {
		log.Fatal("Failed to get absolute path:", err)
	}

	return &Migration{
		databaseURL:    databaseURL,
		migrationsPath: "file://" + strings.ReplaceAll(absPath, "\\", "/"),
	}
}

// Up runs all pending migrations
func (m *Migration) Up() error {
	migrator, err := migrate.New(m.migrationsPath, m.databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("✅ No pending migrations")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✅ All migrations completed successfully")
	return nil
}

// Down rolls back the last migration
func (m *Migration) Down() error {
	migrator, err := migrate.New(m.migrationsPath, m.databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("✅ No migrations to rollback")
			return nil
		}
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Println("✅ Migration rolled back successfully")
	return nil
}

// Status shows the current migration status
func (m *Migration) Status() error {
	migrator, err := migrate.New(m.migrationsPath, m.databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Println("📊 Migration Status: No migrations applied yet")
	} else {
		status := "clean"
		if dirty {
			status = "DIRTY (migration incomplete!)"
		}
		fmt.Printf("📊 Migration Status: Version %d (%s)\n", version, status)
	}

	return nil
}

// Fresh rolls back all migrations and runs them again
func (m *Migration) Fresh() error {
	migrator, err := migrate.New(m.migrationsPath, m.databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Drop(); err != nil {
		return fmt.Errorf("failed to drop all tables: %w", err)
	}

	fmt.Println("✅ All tables dropped successfully")

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✅ Database refreshed and all migrations applied")
	return nil
}

// Seed placeholder untuk future seeding functionality
func (m *Migration) Seed() error {
	fmt.Println("⏳ Seeding functionality coming soon...")
	return nil
}
