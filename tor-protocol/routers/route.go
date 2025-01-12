package routers

import (
	"fmt"

	"tor-protocol/controllers"
	"tor-protocol/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func SetupRoutes(app *fiber.App) {
    api := app.Group("/")
	

    // Middleware for custom headers
    app.Use(middleware.CustomHeaderMiddleware())

    // If the path contains `:port.onion`, 
    // we parse that port and forward the request to the correct server.
    app.All("/:port<[0-9]+>.onion/*", func(c *fiber.Ctx) error {
        // Extract the port (e.g. "9005")
        port := c.Params("port")

        // The remaining path after the `.onion/`
        // For example, if the request is: /9005.onion/foo/bar, `c.Params("*")` will be "foo/bar"
        pathAfterOnion := c.Params("*") // could be empty if just "/9005.onion"
        
        // Construct the target URL that you want to proxy to.
        // E.g. "http://127.0.0.1:9005/foo/bar"
        // If there’s no extra path, this will just be "http://127.0.0.1:9005/"
        target := fmt.Sprintf("http://127.0.0.1:%s/%s", port, pathAfterOnion)

        // Forward the current Fiber context (method, headers, body, etc.) to target
        return proxy.Forward(target)(c)
    })

    // You could also match exactly `:port<[0-9]+>.onion` with no trailing slash:
    app.All("/:port<[0-9]+>.onion", func(c *fiber.Ctx) error {
        port := c.Params("port")
        target := fmt.Sprintf("http://127.0.0.1:%s/", port)
        return proxy.Forward(target)(c)
    })

    // If you still want to serve local `.onion` routes as well:
    app.Get("*.onion", controllers.HomeHandler)

    // Other example routes
    home := api.Group("/home")
    home.Get("/", controllers.ReturnHome)

    // Default route
    send_data := api.Group("/")
    send_data.Get("/", controllers.HomeHandler)
}
