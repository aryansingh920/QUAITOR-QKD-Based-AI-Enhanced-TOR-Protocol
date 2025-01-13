// server.go
package server

import (
	"fmt"
	"log"

	"tor-protocol/client"
	"tor-protocol/config"
	"tor-protocol/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	// "github.com/gofiber/fiber/v2/middleware/limiter"
)

func ServerMain() {
    // Initialize Fiber app
    fmt.Printf("tor-protocol API\n")

    // Example usage of the client, optional
    if err := client.SendRequest(); err != nil {
        log.Printf("Error sending request: %v", err)
    }

    app := fiber.New()

    // Basic middlewares
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "*", // Change to specific domains for better security
        AllowMethods:     "GET,POST,PUT,DELETE",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        ExposeHeaders:    "Content-Length",
        AllowCredentials: true,
    }))
    app.Use(helmet.New())

    // Optionally, you can rate-limit
    // app.Use(limiter.New(limiter.Config{
    //     Max:        10, // Max requests
    //     Expiration: 30 * time.Second, // Time window
    // }))

    // Load environment configuration
    config.LoadConfig()

    // Setup API routes
    routers.SetupRoutes(app)

    // Start server
    port := config.GetPort()
    log.Printf("Server running on port %s", port)
    log.Fatal(app.Listen(":" + port))
}
