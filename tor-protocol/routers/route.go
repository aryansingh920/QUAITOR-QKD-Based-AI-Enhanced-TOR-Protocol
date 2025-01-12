package routers

import (
	"tor-protocol/controllers"
	"tor-protocol/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")

	// Middleware for custom headers
	app.Use(middleware.CustomHeaderMiddleware())

	// Handle `.onion` paths
	app.Get("*.onion", controllers.HomeHandler)

	// User routes
	home := api.Group("/home")
	home.Get("/", controllers.ReturnHome)

	// Default route
	send_data := api.Group("/")
	send_data.Get("/", controllers.HomeHandler)
}
