package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import file source driver
	"gorm.io/gorm"
)

// RunMigrations checks if migrations are applied and applies them if necessary
func RunMigrations(db *gorm.DB) {
	// Get the underlying SQL DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve the underlying SQL DB: %v", err)
	}

	// Create a migration driver for PostgreSQL using the GORM connection
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %v", err)
	}

	// Point to the migrations folder where the migration files are located
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Migration files path
		"postgres", driver)  // Postgres database driver instance

	if err != nil {
		log.Fatalf("Could not initialize migration: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while running migrations: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}
}
