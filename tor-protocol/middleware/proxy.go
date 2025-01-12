package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// ProxyMiddleware handles forwarding requests to the appropriate server based on the `.onion` path.
func ProxyMiddleware(c *fiber.Ctx) error {
	// Extract the port (e.g., "9005")
	port := c.Params("port")

	// The remaining path after `.onion/`
	pathAfterOnion := c.Params("*") // Could be empty if just "/9005.onion"

	// Construct the target URL
	// e.g., "http://127.0.0.1:9005/foo/bar"
	target := fmt.Sprintf("http://127.0.0.1:%s/%s", port, pathAfterOnion)

	// Forward the current Fiber context to the target
	return proxy.Forward(target)(c)
}

// ProxyExactMiddleware handles requests with no trailing slash for `.onion` routes.
func ProxyExactMiddleware(c *fiber.Ctx) error {
	// Extract the port (e.g., "9005")
	port := c.Params("port")

	// Construct the target URL
	// e.g., "http://127.0.0.1:9005/"
	target := fmt.Sprintf("http://127.0.0.1:%s", port)

	// Forward the current Fiber context to the target
	return proxy.Forward(target)(c)
}
