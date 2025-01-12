package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"tor-protocol/client"
	"tor-protocol/config"
	"tor-protocol/routers"
)

func main() {
	// Initialize Fiber app
	fmt.Printf("tor-protocol API\n")
	if err := client.SendRequest(); err != nil {
		log.Printf("Error sending request: %v", err)
	}
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
