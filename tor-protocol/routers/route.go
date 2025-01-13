// routers.go
package routers

import (
	"log"
	"tor-protocol/controllers"
	"tor-protocol/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    api := app.Group("/")

    // Print the path of the request for debugging
    app.Use(func(c *fiber.Ctx) error {
        log.Println("Request URL:", c.OriginalURL())
        return c.Next()
    })

    // Middleware for custom headers
    app.Use(middleware.CustomHeaderMiddleware())

    // Proxy middleware for paths containing `:port.onion`
    app.All("/:port<[0-9]+>.onion/*", middleware.ProxyMiddleware)
    app.All("/:port<[0-9]+>.onion", middleware.ProxyExactMiddleware)

    // Serve local `.onion` routes as well
    app.Get("*.onion", controllers.HomeHandler)

    // Example route group
    home := api.Group("/home")
    home.Get("/", controllers.ReturnHome)

    // Default route
    send_data := api.Group("/")
    send_data.Get("/", controllers.HomeHandler)
}
