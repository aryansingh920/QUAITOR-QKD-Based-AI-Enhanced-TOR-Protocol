package middleware

import (
	"log"
	"tor-protocol/protocol"

	"github.com/gofiber/fiber/v2"
)

// CustomHeaderMiddleware processes custom headers in incoming requests
func CustomHeaderMiddleware() fiber.Handler {
	//log current route with time stamp in format time stamp: route
	// log.Printf(time.Now().Format("2006-01-02 15:04:05") + )

	return func(c *fiber.Ctx) error {
		log.Printf("Current route: %s", c.Route().Path)
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
