package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"tor-protocol/config"
	"tor-protocol/routers"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Load environment configuration
	config.LoadConfig()

	// Setup API routes
	routers.SetupRoutes(app)

	// Start server
	port := config.GetPort()
	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
