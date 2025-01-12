package controllers

import "github.com/gofiber/fiber/v2"

func ReturnHome(c *fiber.Ctx) error {
    if false { // Replace with your actual condition
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "tor-protocol API",
    })
}
