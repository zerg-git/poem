package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	// Remove the command from args so flags parse correctly for subcommands
	os.Args = append(os.Args[:1], os.Args[2:]...)

	switch cmd {
	case "migrate":
		runMigrate()
	case "etl":
		runETL()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: manage <command> [args]")
	fmt.Println("Commands:")
	fmt.Println("  migrate  Run database migrations (users table)")
	fmt.Println("  etl      Run ETL process to import poems (requires chinese-poetry data)")
}

func findProjectRoot() string {
	dir, _ := os.Getwd()

	// 1. Try to find go.mod (Backend root)
	var backendRoot string
	d := dir
	for {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			backendRoot = d
			break
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}

	// 2. If backend root found, check if parent is monorepo root
	if backendRoot != "" {
		parent := filepath.Dir(backendRoot)
		// Check for sibling directories that indicate monorepo root
		if _, err := os.Stat(filepath.Join(parent, "chinese-poetry")); err == nil {
			return parent
		}
		// Or if the folder name is 'backend', assume parent is root
		if filepath.Base(backendRoot) == "backend" {
			return parent
		}
		return backendRoot
	}

	// 3. If go.mod not found, maybe we are at monorepo root?
	if _, err := os.Stat(filepath.Join(dir, "backend")); err == nil {
		return dir
	}

	// Fallback to current directory
	return dir
}

func getDBPath(customPath string) string {
	// If user provided a specific path (not the default flag value), use it
	if customPath != "" && customPath != "poems.db" {
		return customPath
	}

	// Otherwise, default to project root
	root := findProjectRoot()
	return filepath.Join(root, "poems.db")
}
