/*
Updated on 11/01/2025
@author: Aryan

Entry point for the Tor-like simulation. Allows running either a Node or a Client:

	go run main.go <port>            // runs a Node at the given port
	go run main.go <port> client     // runs a Client (for testing requests)
*/
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"tor-protocol/config"
	"tor-protocol/protocol"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage:")
        fmt.Println("  go run main.go <port>           # run a Node")
        fmt.Println("  go run main.go <port> client    # run a Client")
        os.Exit(1)
    }

    // Parse port from the first argument
    portArg := os.Args[1]
    port, err := strconv.Atoi(portArg)
    if err != nil {
        log.Fatalf("Invalid port: %v", err)
    }

    // Check if we have a "client" argument
    var isClient bool
    if len(os.Args) > 2 && os.Args[2] == "client" {
        isClient = true
    }

    // Setup logging to file
    if isClient {
        protocol.SetupLogging(fmt.Sprintf("client_%d.log", port))
    } else {
        protocol.SetupLogging(fmt.Sprintf("node_%d.log", port))
    }

    cfg := config.GetConfig()

    if isClient {
        // Run as client
        log.Printf("Starting Client on port %d...", port)
        protocol.RunClient(port, cfg)
    } else {
        // Run as node
        log.Printf("Starting Node on port %d...", port)
        protocol.RunNode(port, cfg)
    }
}
