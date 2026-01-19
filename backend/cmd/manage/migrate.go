package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/glebarez/sqlite"
)

func runMigrate() {
	// Parse flags
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	dbPath := fs.String("db", "poems.db", "Path to SQLite database")
	fs.Parse(os.Args)

	finalDBPath := getDBPath(*dbPath)

	fmt.Printf("Using database: %s\n", finalDBPath)

	// Connect to database
	db, err := sql.Open("sqlite", finalDBPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Find SQL file
	root := findProjectRoot()
	
	// Search candidates for SQL file
	candidates := []string{
		filepath.Join(root, "backend", "migrations", "create_users_table.sql"),
		filepath.Join(root, "migrations", "create_users_table.sql"),
	}

	var foundPath string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		// Fallback for robustness: try looking around if not found in standard location
		possiblePaths := []string{
			"backend/migrations/create_users_table.sql",
			"migrations/create_users_table.sql",
			"../../migrations/create_users_table.sql",
		}
		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				foundPath = p
				break
			}
		}
		if foundPath == "" {
			log.Fatalf("Failed to read SQL file. Searched in %v and fallback paths", candidates)
		}
	}

	fmt.Printf("Applying migration from: %s\n", foundPath)

	// Read SQL file
	sqlBytes, err := os.ReadFile(foundPath)
	if err != nil {
		log.Fatal("Failed to read SQL file:", err)
	}

	// Execute SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("Failed to execute SQL:", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}
