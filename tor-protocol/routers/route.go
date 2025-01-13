// routers.go
package routers

import (
	"log"

	"tor-protocol/config"

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

    // -----------------------------------------------------------------------
    // 1) Catch-all for "existing route" so that second/third hops keep going.
    // -----------------------------------------------------------------------
    // If X-Tor-Route is already set, we know it's an in-progress multi-hop.
    // We'll skip normal route parsing and forward again via ProxyMiddleware.
    app.All("*", func(c *fiber.Ctx) error {
        existingRoute := c.Get(config.CustomHeaderKey)
        if existingRoute != "" {
            return middleware.ProxyMiddleware(c)
        }
        return c.Next()
    })

    // Middleware for custom headers (optional usage)
    app.Use(middleware.CustomHeaderMiddleware())

    // -----------------------------------------------------------------------
    // 2) These routes handle the "first hop" – building or starting the route
    //    if no X-Tor-Route exists yet.
    // -----------------------------------------------------------------------
    // Proxy middleware for paths containing `:port.onion`
    app.All("/:port<[0-9]+>.onion/*", middleware.ProxyMiddleware)
    app.All("/:port<[0-9]+>.onion", middleware.ProxyExactMiddleware)

    app.All("/:port<[0-9]+>/*", middleware.ProxyMiddleware)
    app.All("/:port<[0-9]+>", middleware.ProxyExactMiddleware)

    // Serve local `.onion` routes as well
    app.Get("*.onion", controllers.HomeHandler)
    app.Get("*", controllers.HomeHandler)

    // Example route group
    home := api.Group("/home")
    home.Get("/", controllers.ReturnHome)

    // Default route
    send_data := api.Group("/")
    send_data.Get("/", controllers.HomeHandler)
}
