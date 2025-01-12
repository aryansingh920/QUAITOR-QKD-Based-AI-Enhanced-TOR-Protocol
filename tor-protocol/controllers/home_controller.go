package controllers

import (
	"fmt"
	"tor-protocol/utils/protocol"

	"github.com/gofiber/fiber/v2"
)

func ReturnHome(c *fiber.Ctx) error {

	customHeader, ok := c.Locals("customHeader").(*protocol.CustomHeader)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	// Process the request using the custom header
	fmt.Printf("Received request with RouteID: %d, PayloadSize: %d\n", customHeader.RouteID, customHeader.PayloadSize)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "tor-protocol API",
	})
}
