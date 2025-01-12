package middleware

import (
	"tor-protocol/protocol"

	"github.com/gofiber/fiber/v2"
)

// CustomHeaderMiddleware processes custom headers in incoming requests
func CustomHeaderMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("X-QUAITOR-Protocol")
		if header == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing custom protocol header")
		}

		// Decode header
		data := []byte(header)
		customHeader, err := protocol.Deserialize(data)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid custom protocol header")
		}

		// Add header to context for downstream handlers
		c.Locals("customHeader", customHeader)

		return c.Next()
	}
}
