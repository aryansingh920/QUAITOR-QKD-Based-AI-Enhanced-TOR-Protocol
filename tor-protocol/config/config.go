package config

import (
	"log"
	"os"
)

func LoadConfig() {
	// Load environment variables if .env file exists
	if err := os.Setenv("LOG_LEVEL", "info"); err != nil {
		log.Println("Environment variables loaded with default configurations")
	}
}

func GetPort() string {
	// Check if a command-line argument is passed for the port
	args := os.Args
	if len(args) > 1 {
		port := args[1]
		log.Printf("Using port from command-line argument: %s\n", port)
		return port
	}

	// Fallback to default port if no argument is provided
	defaultPort := "3000"
	log.Printf("No command-line argument provided, using default port: %s\n", defaultPort)
	return defaultPort
}
