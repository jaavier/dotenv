// Command example demonstrates how to load and read environment variables
// with the dotenv package.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jaavier/dotenv"
)

func main() {
	// Load the default .env file. Existing environment variables are NOT
	// overridden — the file only fills in what is missing (12-factor safe).
	if err := dotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	// Using GetOrPanic for required variables.
	apiKey := dotenv.GetOrPanic("API_KEY")
	fmt.Printf("API_KEY loaded: %s...\n", apiKey[:min(10, len(apiKey))])

	// Using GetOrDefault for optional variables.
	port := dotenv.GetOrDefault("PORT", "8080")
	fmt.Printf("Server will run on port: %s\n", port)

	// Using Get (same as os.Getenv).
	if dotenv.Get("DEBUG") == "true" {
		fmt.Println("Debug mode enabled")
	}

	// Parse without side effects: inspect values without touching the
	// process environment. Great for tests and validation.
	sample := "FEATURE_X=on\nFEATURE_Y=off # disabled for now"
	vars, err := dotenv.Parse(strings.NewReader(sample))
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}
	fmt.Printf("Parsed (no env mutation): %v\n", vars)

	// Overload when the file should be the source of truth (opt-in override).
	if err := dotenv.Overload("config/.env.production"); err != nil {
		log.Printf("Note: production config not found, using defaults\n")
	}

	// Example: build database configuration.
	dbConfig := map[string]string{
		"host": dotenv.GetOrPanic("DB_HOST"),
		"port": dotenv.GetOrDefault("DB_PORT", "5432"),
		"name": dotenv.GetOrPanic("DB_NAME"),
		"user": dotenv.GetOrPanic("DB_USER"),
	}

	fmt.Printf("Database configured: %s@%s:%s/%s\n",
		dbConfig["user"], dbConfig["host"], dbConfig["port"], dbConfig["name"])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
