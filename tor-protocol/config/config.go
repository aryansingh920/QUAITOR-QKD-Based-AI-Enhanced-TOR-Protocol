// config.go
package config

import (
	"log"
	"os"
)

const (
	PortStart = 8801 
	PortEnd   = 8820 
    DefaultLink = "http://127.0.0.1"
    RandomDelayUpperLimit = 5000
    CustomHeaderKey = "X-Tor-Route"
)


func LoadConfig() {
    // Load environment variables or defaults
    // You can extend this to parse .env files with e.g. github.com/joho/godotenv if desired
    if err := os.Setenv("LOG_LEVEL", "info"); err != nil {
        log.Println("Environment variables loaded with default configurations")
    }
}

// GetPort returns the server port from command-line arg if provided,
// otherwise it defaults to "3000".
func GetPort() string {
    args := os.Args
    if len(args) > 1 {
        port := args[1]
        log.Printf("Using port from command-line argument: %s\n", port)
        return port
    }

    defaultPort := "3000"
    log.Printf("No command-line argument provided, using default port: %s\n", defaultPort)
    return defaultPort
}
