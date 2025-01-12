package routers

import (
	"tor-protocol/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")

	// User routes
	home := api.Group("/")
	home.Get("/", controllers.ReturnHome)

	// Product routes with middleware
	// product := api.Group("/product", middleware.AuthMiddleware)
	// product.Get("/", controllers.GetProducts)
}
