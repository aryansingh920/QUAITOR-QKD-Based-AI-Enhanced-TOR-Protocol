/*
Updated on 11/01/2025
@author: Aryan

Implements the Node logic:
  - Listening on a TCP port
  - Accepting connections and forwarding data to the next node
  - Potentially generating random (dummy) traffic
*/
package protocol

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
	"tor-protocol/config"
)

type NodeRole string

const (
    EntryNode NodeRole = "ENTRY"
    RelayNode NodeRole = "RELAY"
    ExitNode  NodeRole = "EXIT"
)

// Node represents a Tor-like node
type Node struct {
    ID             string
    Role           NodeRole
    Port           int
    listener       net.Listener
    stopCh         chan struct{}
    wg             sync.WaitGroup
    ephemeralPub   []byte
    ephemeralPriv  []byte

    // Just an example of holding "service" data
    // "9000.onion" -> "Hello from port 9000"
    HtmlContent string
}

// Start begins listening on n.Port and handles incoming traffic
func (n *Node) Start() error {
    var err error
    addr := fmt.Sprintf(":%d", n.Port)

    n.listener, err = net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("Node %s failed to listen on port %d: %v", n.ID, n.Port, err)
    }

    // Generate ephemeral key pair for demonstration
    n.ephemeralPub, n.ephemeralPriv, err = GenerateEphemeralKeyPair()
    if err != nil {
        log.Printf("[%s] Failed to generate ephemeral keys: %v\n", n.ID, err)
    } else {
        log.Printf("[%s] Ephemeral keys generated. PublicKey length=%d, PrivateKey length=%d\n",
            n.ID, len(n.ephemeralPub), len(n.ephemeralPriv))
    }

    // Log role & listening
    log.Printf("[%s] %s node listening on %s\n", n.ID, n.Role, addr)

    // Accept connections
    n.wg.Add(1)
    go func() {
        defer n.wg.Done()
        for {
            conn, err := n.listener.Accept()
            if err != nil {
                select {
                case <-n.stopCh:
                    // shutting down
                    return
                default:
                    log.Printf("[%s] Accept error: %v", n.ID, err)
                    continue
                }
            }

            // handle in a new goroutine
            n.wg.Add(1)
            go func(c net.Conn) {
                defer n.wg.Done()
                n.handleConnection(c)
            }(conn)
        }
    }()
    return nil
}

// handleConnection processes incoming RelayCells
func (n *Node) handleConnection(conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 4096)
    for {
        // set a read deadline to avoid idle connections
        conn.SetDeadline(time.Now().Add(300 * time.Second))
        bytesRead, err := conn.Read(buffer)
        if err != nil {
            // could log or ignore
            return
        }
        data := buffer[:bytesRead]

        relayCell, err := ParseRelayCell(data)
        if err != nil {
            log.Printf("[%s] Error parsing RelayCell: %v\n", n.ID, err)
            return
        }

        // For demonstration, skip real decryption
        // In real Tor, we'd decrypt a layer here with ephemeralPriv

        // If I'm the ExitNode for the requested onion service, handle the request
        if n.Role == ExitNode && relayCell.IsExitRequest {
            // parse the "onion" request from the payload
            onionAddress := string(relayCell.Payload)
            log.Printf("[%s] EXIT node handling request for '%s'\n", n.ID, onionAddress)
            // "Serve" the HTML content
            // Example: if the onionAddress is "9000.onion", respond with n.HtmlContent
            response := []byte(n.HtmlContent)
            // Send the response back in the same RelayCell (in real Tor you'd wrap it)
            relayCell.Payload = response
            relayCell.IsExitResponse = true
            relayCell.IsExitRequest = false

            // Forward back to the previous node (which must have been set in relayCell.PrevAddr)
            // Notice we do NOT read the entire chain; we only know the previous address to send back
            // If there's no PrevAddr, that means we can't respond
            if relayCell.PrevAddr != "" {
                err = forwardToAddress(relayCell, relayCell.PrevAddr)
                if err != nil {
                    log.Printf("[%s] Error sending exit response back: %v\n", n.ID, err)
                }
            }
            continue
        }

        // If it's an ExitResponse, it means data is returning to the client, so we forward it back
        if relayCell.IsExitResponse {
            if relayCell.NextAddr != "" {
                // forward the data back up the chain
                err = forwardToAddress(relayCell, relayCell.NextAddr)
                if err != nil {
                    log.Printf("[%s] Error forwarding exit response: %v\n", n.ID, err)
                }
            }
            continue
        }

        // Otherwise, I'm not the exit node, so just forward to the next hop (if any)
        if relayCell.NextAddr != "" {
            // log the forwarding action
            log.Printf("[%s] Forwarding data to next node: %s\n", n.ID, relayCell.NextAddr)
            err = forwardToAddress(relayCell, relayCell.NextAddr)
            if err != nil {
                log.Printf("[%s] Error forwarding to next node: %v\n", n.ID, err)
            }
        }
    }
}

// Stop closes the node's listener and waits for goroutines to finish
func (n *Node) Stop() {
    close(n.stopCh)
    if n.listener != nil {
        _ = n.listener.Close()
    }
    n.wg.Wait()
    log.Printf("[%s] Node shut down.", n.ID)
}

// GenerateRandomTraffic simulates dummy traffic to random known nodes at intervals
func (n *Node) GenerateRandomTraffic(intervalSeconds int, knownPorts []int) {
    ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // pick a random port for dummy traffic
            targetPort := knownPorts[rand.Intn(len(knownPorts))]
            // skip sending to self
            if targetPort == n.Port {
                continue
            }

            // Construct a dummy RelayCell
            cell := &RelayCell{
                PrevAddr: "", // no previous if I'm the origin
                NextAddr: fmt.Sprintf("127.0.0.1:%d", targetPort),
                Payload:  []byte("[DUMMY_TRAFFIC] Hello from " + n.ID),
            }

            // Attempt to dial
            err := forwardToAddress(cell, cell.NextAddr)
            if err != nil {
                log.Printf("[%s] Dummy traffic error: %v\n", n.ID, err)
            } else {
                log.Printf("[%s] Sent dummy traffic to port %d\n", n.ID, targetPort)
            }

        case <-n.stopCh:
            return
        }
    }
}



// RunNode creates a single Node with (possibly) random or default role, starts it,
// and optionally starts random traffic if enabled in the config.
func RunNode(port int, cfg *config.Config) {
    // Create a Node. For demonstration, we'll guess a random role if not forced by an ENV var.
    nodeRole := getEnvOrRandomRole()
    node := &Node{
        ID:   fmt.Sprintf("Node%d", port),
        Port: port,
        Role: nodeRole,
        stopCh: make(chan struct{}),
        // Example HTML content
        HtmlContent: fmt.Sprintf("<html><body><h1>Hello from port %d (role=%s)</h1></body></html>", port, nodeRole),
    }

    err := node.Start()
    if err != nil {
        log.Fatalf("Node %s encountered an error: %v", node.ID, err)
    }

    // If random traffic is enabled, start generating it
    if cfg.EnableRandomTraffic {
        log.Printf("[RunNode] Starting random traffic in node %s...", node.ID)
        go node.GenerateRandomTraffic(cfg.RandomTrafficInterval, cfg.KnownPorts)
    }

    // Keep running unless forcibly stopped. In a real program, you might use signals or a user prompt.
    for {
        time.Sleep(1 * time.Hour)
    }
}

// getEnvOrRandomRole checks if TOR_NODE_ROLE is set, otherwise picks a random role
func getEnvOrRandomRole() NodeRole {
    envRole := getEnv("TOR_NODE_ROLE", "")
    if envRole != "" {
        switch envRole {
        case "ENTRY":
            return EntryNode
        case "EXIT":
            return ExitNode
        default:
            return RelayNode
        }
    }
    // If not set, randomly choose
    roles := []NodeRole{EntryNode, RelayNode, ExitNode}
    return roles[rand.Intn(len(roles))]
}

func getEnv(key, fallback string) string {
    val := ""
    if v := getenv(key); v != "" {
        val = v
    } else {
        val = fallback
    }
    return val
}
func getenv(key string) string {
    return "" // you can do: return os.Getenv(key)
}
