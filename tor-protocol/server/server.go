// server.go
package server

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"tor-protocol/client"
	"tor-protocol/config"
	"tor-protocol/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	// "github.com/gofiber/fiber/v2/middleware/limiter"
)

func ServerMain() {
	// Load environment configuration
	config.LoadConfig()

	// Get current port
	port := config.GetPort()

	// Set up logging to both terminal and file
	logFilePath := filepath.Join("logs", fmt.Sprintf("server-%s.log", port))
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", logFilePath, err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("tor-protocol API started on port %s\n", port)

	// Initialize Fiber app
	app := fiber.New()

	// Basic middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // For production, change to specific domains
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))
	app.Use(helmet.New())

	// Optionally, you can rate-limit
	// app.Use(limiter.New(limiter.Config{
	//     Max:        10,                // Max requests
	//     Expiration: 30 * time.Second,  // Time window
	// }))

	// Example usage of the client, optional
	if err := client.SendRequest(); err != nil {
		log.Printf("Error sending request: %v", err)
	}

	// Setup API routes
	routers.SetupRoutes(app)

	// Start server
	log.Printf("Server running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
}
