/*
Updated on 11/01/2025
@author: Aryan

Provides configuration (ports, random traffic, etc.).
*/
package config

import (
	"os"
	"strconv"
)

// Config holds basic config info. You can expand this to read from environment or files.
type Config struct {
    // Default (known) ports in the network. In real usage, you'd discover them dynamically.
    KnownPorts []int

    // If set, random cover traffic will be generated
    EnableRandomTraffic bool

    // Interval for generating random traffic (seconds)
    RandomTrafficInterval int
}

// GetConfig reads config from environment or uses defaults
func GetConfig() *Config {
    // Example known ports. Adjust as needed or use environment.
    defaultKnownPorts := []int{9000, 9001, 9002, 9003, 9004, 9005}

    // Enable random traffic from env
    enableRandomTraffic := false
    if os.Getenv("TOR_LIKE_ENABLE_RANDOM_TRAFFIC") == "true" {
        enableRandomTraffic = true
    }

    // Random traffic interval from env
    randomTrafficInterval := 10
    if intervalStr := os.Getenv("TOR_LIKE_RANDOM_TRAFFIC_INTERVAL"); intervalStr != "" {
        if val, err := strconv.Atoi(intervalStr); err == nil && val > 0 {
            randomTrafficInterval = val
        }
    }

    return &Config{
        KnownPorts:           defaultKnownPorts,
        EnableRandomTraffic:  enableRandomTraffic,
        RandomTrafficInterval: randomTrafficInterval,
    }
}
