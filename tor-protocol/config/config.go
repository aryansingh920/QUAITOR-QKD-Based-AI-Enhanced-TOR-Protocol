package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}
	return port
}
