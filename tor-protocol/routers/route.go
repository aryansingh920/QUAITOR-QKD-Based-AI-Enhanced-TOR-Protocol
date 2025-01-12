package routers

import (
	"tor-protocol/controllers"
	"tor-protocol/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/")

	app.Use(middleware.CustomHeaderMiddleware())

	// User routes
	home := api.Group("/")
	home.Get("/", controllers.ReturnHome)

	send_data := api.Group("/send_data")
	send_data.Get("/", controllers.HomeHandler)
	

	// Product routes with middleware
	// product := api.Group("/product", middleware.AuthMiddleware)
	// product.Get("/", controllers.GetProducts)
}


