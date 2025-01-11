/*
Created on 11/01/2025

@author: Aryan

Filename: config.go

Relative Path: tor-protocol/config/config.go
*/

package config

import (
	"os"
	"strconv"
)

type Config struct {
    Ports              []int
    RandomPathLength   int
    EnableRandomTraffic bool
    RandomTrafficInterval int // in seconds
}

// GetConfig reads some basic configuration from environment variables or provides defaults.
func GetConfig() *Config {
    // Example default ports
    defaultPorts := []int{9001, 9002, 9003, 9004, 9005,9006,9007,9008,9009,9010}

    // Optional: read from ENV, or just keep default
    strPathLength := os.Getenv("TOR_LIKE_PATH_LENGTH")
    pathLength, err := strconv.Atoi(strPathLength)
    if err != nil {
        pathLength = 0 // 0 indicates random
    }

    // Enable or disable random traffic
    enableRandomTrafficStr := os.Getenv("TOR_LIKE_ENABLE_RANDOM_TRAFFIC")
    enableRandomTraffic := false
    if enableRandomTrafficStr == "1" || enableRandomTrafficStr == "true" {
        enableRandomTraffic = true
    }

    // Random traffic interval
    randomTrafficIntervalStr := os.Getenv("TOR_LIKE_RANDOM_TRAFFIC_INTERVAL")
    randomTrafficInterval, err := strconv.Atoi(randomTrafficIntervalStr)
    if err != nil || randomTrafficInterval <= 0 {
        randomTrafficInterval = 10 // default to 10 seconds if not set properly
    }

    return &Config{
        Ports:                 defaultPorts,
        RandomPathLength:      pathLength,
        EnableRandomTraffic:   enableRandomTraffic,
        RandomTrafficInterval: randomTrafficInterval,
    }
}
