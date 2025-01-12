package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ProxyOnionRequests handles requests to *.onion and forwards them dynamically
func ProxyOnionRequests(c *fiber.Ctx) error {
	// Extract port from the .onion path
	onionPath := c.Params("*") // Extracts the dynamic part of *.onion
	parts := strings.Split(onionPath, ".")
	if len(parts) != 2 || parts[1] != "onion" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid .onion path")
	}

	// Parse the target port
	targetPort, err := strconv.Atoi(parts[0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid port in .onion path")
	}

	// Construct target URL
	targetURL := fmt.Sprintf("http://127.0.0.1:%d/", targetPort)
	log.Printf("Proxying request to %s", targetURL)

	// Forward request to target server
	resp, err := http.Get(targetURL) // Sends a GET request to the target server
	if err != nil {
		log.Printf("Error proxying to %s: %v", targetURL, err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to reach target server")
	}
	defer resp.Body.Close()

	// Read and return the response from the target server
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from %s: %v", targetURL, err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading target response")
	}

	// Set headers from the target response (optional)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Return the response to the original client
	return c.Status(resp.StatusCode).Send(body)
}
