// config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	PortStart = 8801 
	PortEnd  = 8805
    DefaultLink = "http://127.0.0.1"
    RandomDelayUpperLimit = 5000
    CustomHeaderKey = "X-Tor-Route"
)




func LoadConfig() {
    var err error
    // Load environment variables or defaults
    // You can extend this to parse .env files with e.g. github.com/joho/godotenv if desired


    PortStartEnv, err := getEnvAsInt("start_port", 8801)
	if err != nil {
		log.Printf("Error parsing PORT_START, using default 8801: %v\n", err)
	}

	PortStart = PortStartEnv


	PortEndEnv, err := getEnvAsInt("end_port", 8820)
	if err != nil {
		log.Printf("Error parsing PORT_END, using default 8820: %v\n", err)
	}

	PortEnd = PortEndEnv 

    fmt.Printf("At Config: PortStart: %d, PortEnd: %d\n", PortStart, PortEnd)
    log.Printf("At Config: PortStart: %d, PortEnd: %d\n", PortStart, PortEnd)


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

// getEnvAsInt retrieves the value of the environment variable as an integer or returns a default value.
func getEnvAsInt(key string, defaultValue int) (int, error) {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(valueStr)
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
