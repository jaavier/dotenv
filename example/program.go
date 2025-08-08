package main

import (
	"fmt"
	"log"

	"github.com/jaavier/dotenv"
)

func main() {
	if err := dotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v\n", err)
	}

	// Using GetOrPanic for required variables
	apiKey := dotenv.GetOrPanic("API_KEY")
	fmt.Printf("API_KEY loaded: %s...\n", apiKey[:min(10, len(apiKey))])

	// Using GetOrDefault for optional variables
	port := dotenv.GetOrDefault("PORT", "8080")
	fmt.Printf("Server will run on port: %s\n", port)

	// Using Get (same as os.Getenv)
	debugMode := dotenv.Get("DEBUG")
	if debugMode == "true" {
		fmt.Println("Debug mode enabled")
	}

	// Load additional config file with options
	opts := &dotenv.Options{
		Override: false,
		Required: false,
	}
	if err := dotenv.LoadWithOptions(opts, "config/.env.production"); err != nil {
		log.Printf("Note: Production config not found, using defaults\n")
	}

	// Example: Get database configuration with panic for required fields
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