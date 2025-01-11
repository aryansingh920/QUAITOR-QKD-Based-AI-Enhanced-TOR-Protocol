/*
Created on 11/01/2025

@author: Aryan

Filename: main.go

Relative Path: tor-protocol/main.go
*/
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"tor-protocol/config"
	"tor-protocol/protocol"
)

func main() {
    // Read config
    cfg := config.GetConfig()
	// fmt.Printf("Config: %+v\n", cfg)

    // Seed random
    rand.Seed(time.Now().UnixNano())

    // Prepare the nodes
    nodes := protocol.BootstrapNodes(cfg)

    // Start each node as a goroutine
    for i := range nodes {
        go func(n *protocol.Node) {
            if err := n.Start(); err != nil {
                log.Fatalf("Node %s encountered an error: %v", n.ID, err)
            }
        }(nodes[i])
    }

    // Give the nodes a moment to start up
    time.Sleep(2 * time.Second)

    // OPTIONAL: Start generating random/dummy traffic among nodes (if enabled in config)
    if cfg.EnableRandomTraffic {
        log.Println("Starting random traffic generation among nodes...")
        for i := range nodes {
            go nodes[i].GenerateRandomTraffic(cfg.RandomTrafficInterval)
        }
    }

    // Build a random path of length X (from config, or random 2-5, e.g.)
    pathLength := cfg.RandomPathLength
    if pathLength <= 0 {
        // fallback to random between 3 and len(nodes) if not specified
        pathLength = rand.Intn(len(nodes)-2) + 2
    }

    // Build circuit: pick random nodes for the path
    circuit, err := protocol.BuildCircuit(nodes, pathLength)
    if err != nil {
        log.Fatalf("Error building circuit: %v", err)
    }

    // Display the chosen path
    fmt.Printf("Circuit path: ")
    for _, node := range circuit {
        fmt.Printf("%s -> ", node.ID)
    }
    fmt.Println()

    // Send a test message through the circuit
    testMessage := "Hello from our Tor-like client!"
    err = protocol.SendThroughCircuit(circuit, testMessage)
    if err != nil {
        log.Fatalf("Error sending message through circuit: %v", err)
    }

    // Keep the program running for demonstration. In a real scenario,
    // you'd let the nodes continue listening or gracefully shut down.
    time.Sleep(5 * time.Second)
    // fmt.Println("Shutting down...")

    // Stop all nodes
    // for _, n := range nodes {
    //     n.Stop()
    // }
}
