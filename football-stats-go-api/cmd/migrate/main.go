package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		migrationsPath = flag.String("path", "file://migrations", "Path to migrations folder")
		dbURL          = flag.String("db", "", "Database URL (postgres://user:pass@host:port/dbname?sslmode=disable)")
		direction      = flag.String("direction", "up", "Migration direction: up or down")
	)
	flag.Parse()

	if *dbURL == "" {
		*dbURL = os.Getenv("DATABASE_URL")
	}

	if *dbURL == "" {
		log.Fatal("Database URL is required. Use -db flag or DATABASE_URL env var")
	}

	m, err := migrate.New(*migrationsPath, *dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", *direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get version: %v", err)
	}

	fmt.Printf("✓ Migration completed successfully\n")
	fmt.Printf("  Current version: %d\n", version)
	fmt.Printf("  Dirty: %v\n", dirty)
}
