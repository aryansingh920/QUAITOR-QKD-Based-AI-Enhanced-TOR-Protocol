/*
Updated on 11/01/2025
@author: Aryan

Implements the Client logic:
  - Building a random circuit
  - Sending a request (like "GET 9000.onion") through the circuit
  - Receiving the response
*/
package protocol

import (
	"fmt"
	"log"
	"net"
	"time"

	"tor-protocol/config"
)

// RunClient is called from main if we're in 'client' mode
func RunClient(port int, cfg *config.Config) {
    log.Printf("Client on port %d starting up...", port)

    // Hardcode a test onion address to request
    targetOnion := fmt.Sprintf("%d.onion", cfg.KnownPorts[0])
    log.Printf("[Client] Attempting to request HTML from: %s", targetOnion)

    // Build a random circuit of length 3 (or whatever you'd like)
    circuitLength := 3
    if circuitLength > len(cfg.KnownPorts) {
        circuitLength = len(cfg.KnownPorts)
    }
    nodes, err := BuildCircuit(cfg.KnownPorts, circuitLength)
    if err != nil {
        log.Fatalf("[Client] Failed to build circuit: %v", err)
    }

    // Start each node (in a real scenario, the nodes are already up 
    // and we wouldn't start them from the client).
    // For demonstration, we'll forcibly start them so we can test quickly.
    for _, n := range nodes {
        go func(nd *Node) {
            err := nd.Start()
            if err != nil {
                log.Fatalf("Node %s encountered an error: %v", nd.ID, err)
            }
        }(n)
    }

    // Give them a second to start
    time.Sleep(1 * time.Second)

    // Optional: launch random traffic if enabled
    if cfg.EnableRandomTraffic {
        log.Println("[Client] Starting random traffic in all nodes for cover traffic.")
        for _, n := range nodes {
            go n.GenerateRandomTraffic(cfg.RandomTrafficInterval, cfg.KnownPorts)
        }
    }

    // Actually send the request through the circuit
    err = sendRequestThroughCircuit(nodes, targetOnion)
    if err != nil {
        log.Printf("[Client] Error sending request: %v", err)
        return
    }

    // Wait a bit for the response to come back, for demonstration
    time.Sleep(3 * time.Second)

    // Stop nodes
    for _, n := range nodes {
        n.Stop()
    }
    log.Printf("[Client] All done. Exiting.")
}

// sendRequestThroughCircuit simulates an onion request
func sendRequestThroughCircuit(circuit []*Node, onionAddress string) error {
    if len(circuit) == 0 {
        return fmt.Errorf("circuit is empty")
    }

    // We set up consecutive addresses. For node i, next is node i+1
    // The first node has no PrevAddr (origin = client).
    // The last node has no NextAddr (it is the exit).
    // We also mark the final cell as IsExitRequest.

    for i := 0; i < len(circuit)-1; i++ {
        circuit[i].Role = EntryNode
        circuit[i+1].Role = ExitNode
        circuit[i].Role = RelayNode
    }
    // Actually, let's revert to the normal approach:
    circuit[0].Role = EntryNode
    circuit[len(circuit)-1].Role = ExitNode
    for i := 1; i < len(circuit)-1; i++ {
        circuit[i].Role = RelayNode
    }

    // Build the chain of NextAddr
    for i := 0; i < len(circuit)-1; i++ {
        // nextNodePort := circuit[i+1].Port
        circuit[i].HtmlContent = fmt.Sprintf("<html><body><h1>Hello from port %d</h1></body></html>", circuit[i].Port)
        circuit[i].ID = fmt.Sprintf("Node%d", circuit[i].Port)
        // circuit[i].Role = circuit[i].Role // assigned above
        // We don't actually store the next address in the node itself; 
        // we do it in the RelayCell payload at runtime.

        // The exit node gets its own HTML content
        if i == len(circuit)-2 {
            circuit[i+1].HtmlContent = fmt.Sprintf("<html><body><h1>Hello from port %d</h1></body></html>", circuit[i+1].Port)
            circuit[i+1].ID = fmt.Sprintf("Node%d", circuit[i+1].Port)
        }

        // Start them if not already started. (In real usage, they'd already be running.)
    }

    // Start with the "entry node"
    entryNode := circuit[0]
    entryAddr := fmt.Sprintf("127.0.0.1:%d", entryNode.Port)
    conn, err := net.DialTimeout("tcp", entryAddr, 2*time.Second)
    if err != nil {
        return fmt.Errorf("unable to dial entry node at %s: %v", entryAddr, err)
    }
    defer conn.Close()

    // Construct the RelayCell
    // The NextAddr is the second node, or if only 1 node, it's the exit (same node).
    // We'll just encode the entire chain step by step. But to be minimal, we only set NextAddr to circuit[1].
    var nextAddr string
    if len(circuit) > 1 {
        nextAddr = fmt.Sprintf("127.0.0.1:%d", circuit[1].Port)
    } else {
        // single-node circuit (not typical in Tor, but let's allow)
        nextAddr = ""
    }

    requestCell := &RelayCell{
        PrevAddr:      "", // client has no address
        NextAddr:      nextAddr,
        Payload:       []byte(onionAddress), // the onion address we want
        IsExitRequest: (len(circuit) == 1),   // if single node, it's also exit
    }

    // Write the request to the entry node
    data, err := requestCell.Serialize()
    if err != nil {
        return fmt.Errorf("serialize request: %v", err)
    }
    _, err = conn.Write(data)
    if err != nil {
        return fmt.Errorf("writing to entry node: %v", err)
    }
    log.Printf("[Client] Sent onion request '%s' to entry node %s", onionAddress, entryNode.ID)

    return nil
}
