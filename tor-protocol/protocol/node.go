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
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
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
    conn.SetDeadline(time.Now().Add(300 * time.Second))
    bytesRead, err := conn.Read(buffer)
    if err != nil {
        return
    }
    data := buffer[:bytesRead]

    // NEW: Check if it's an HTTP GET
    if isHTTPRequest(data) {
        onionAddr := parseOnionPath(data)
        if onionAddr == "" {
            // If we can't parse the onion from the path, return a 400
            conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\nMissing onion path\n"))
            return
        }
        log.Printf("[%s] Got HTTP request for onion: %s\n", n.ID, onionAddr)

        // Build circuit + fetch
        htmlResp, err := n.fetchOnionViaCircuit(onionAddr)
        if err != nil {
            log.Printf("[%s] Error fetching onion: %v\n", n.ID, err)
            conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\nFailed to fetch onion\n"))
            return
        }

        // Write back a simple HTTP 200
        responseHeaders := "HTTP/1.1 200 OK\r\n" +
            "Content-Type: text/html\r\n" +
            "Connection: close\r\n\r\n"
        conn.Write([]byte(responseHeaders))
        conn.Write(htmlResp)
        return
    }

    // OLD logic (existing RelayCell handling) remains unchanged below:
    relayCell, err := ParseRelayCell(data)
    if err != nil {
        log.Printf("[%s] Error parsing RelayCell: %v\n", n.ID, err)
        return
    }
    fmt.Printf("%+v\n", relayCell)
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
func getenv(_ string) string {
    return "" // you can do: return os.Getenv(key)
}

// isHTTPRequest checks if the incoming data starts with "GET "
func isHTTPRequest(data []byte) bool {
    return len(data) > 4 && strings.HasPrefix(string(data), "GET ")
}

// parseOnionPath extracts "9005.onion" from "GET /9005.onion HTTP/1.1"
func parseOnionPath(data []byte) string {
    // Read only the first line: "GET /9005.onion HTTP/1.1"
    scanner := bufio.NewScanner(strings.NewReader(string(data)))
    if !scanner.Scan() {
        return ""
    }
    line := scanner.Text() // e.g. "GET /9005.onion HTTP/1.1"

    parts := strings.Split(line, " ")
    // parts[0] = "GET", parts[1] = "/9005.onion", parts[2] = "HTTP/1.1"
    if len(parts) < 2 {
        return ""
    }
    path := parts[1] // "/9005.onion"
    path = strings.TrimPrefix(path, "/") // remove leading slash -> "9005.onion"
    return path
}

// fetchOnionViaCircuit is a simplified snippet that
// builds a circuit, sends a request, and returns the HTML response.
func (n *Node) fetchOnionViaCircuit(onion string) ([]byte, error) {
    known := []int{9000, 9001, 9002, 9003, 9004, 9005}
    circuitLength := 3
    if circuitLength > len(known) {
        circuitLength = len(known)
    }

    nodes, err := BuildCircuit(known, circuitLength)
    if err != nil {
        return nil, fmt.Errorf("build circuit: %v", err)
    }

    // Start nodes and handle errors
    for _, nd := range nodes {
        go func(node *Node) {
            if err := node.Start(); err != nil {
                log.Printf("Failed to start node %s: %v", node.ID, err)
            }
        }(nd)
    }

    // Wait briefly for nodes to start
    time.Sleep(1 * time.Second)

    // Dial entry node
    entryNode := nodes[0]
    entryAddr := fmt.Sprintf("127.0.0.1:%d", entryNode.Port)
    conn, err := net.DialTimeout("tcp", entryAddr, 2*time.Second)
    if err != nil {
        // Stop all nodes on failure
        for _, nd := range nodes {
            nd.Stop()
        }
        return nil, fmt.Errorf("dial entry node: %v", err)
    }
    defer conn.Close()

    // Set up the NextAddr for the entry node
    var nextAddr string
    if len(nodes) > 1 {
        nextAddr = fmt.Sprintf("127.0.0.1:%d", nodes[1].Port)
    }

    // Prepare the request cell
    requestCell := &RelayCell{
        PrevAddr:      "",
        NextAddr:      nextAddr,
        Payload:       []byte(onion),
        IsExitRequest: (len(nodes) == 1),
    }

    // Serialize and send the request
    data, err := requestCell.Serialize()
    if err != nil {
        return nil, fmt.Errorf("serialize: %v", err)
    }
    if _, err := conn.Write(data); err != nil {
        return nil, fmt.Errorf("write entry node: %v", err)
    }

    // Read the response
    buf := make([]byte, 4096)
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    nBytes, err := conn.Read(buf)
    if err != nil {
        return nil, fmt.Errorf("read from entry node: %v", err)
    }
    respCell, err := ParseRelayCell(buf[:nBytes])
    if err != nil {
        return nil, fmt.Errorf("parse relay cell: %v", err)
    }
    if !respCell.IsExitResponse {
        return nil, fmt.Errorf("did not receive exit response cell")
    }

    // Stop all nodes after the request
    for _, nd := range nodes {
        nd.Stop()
    }

    return respCell.Payload, nil
}

