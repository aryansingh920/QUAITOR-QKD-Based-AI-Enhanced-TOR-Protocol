package protocol

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// NodeRole represents the role of a Node: entry, middle (relay), or exit.
type NodeRole string

const (
    EntryNode NodeRole = "ENTRY"
    RelayNode NodeRole = "RELAY"
    ExitNode  NodeRole = "EXIT"
)

// Node represents a Tor-like node
type Node struct {
    ID         string
    Role       NodeRole
    Address    string
    Port       int
    listener   net.Listener
    stopCh     chan struct{}
    wg         sync.WaitGroup
}

// Start begins listening on the specified port
func (n *Node) Start() error {
    var err error
    addr := fmt.Sprintf(":%d", n.Port)

    n.listener, err = net.Listen("tcp", addr)
    if err != nil {
        return fmt.Errorf("Node %s failed to listen on port %d: %v", n.ID, n.Port, err)
    }

    log.Printf("[%s] %s node listening on %s\n", n.ID, n.Role, addr)

    // accept connections
    n.wg.Add(1)
    go func() {
        defer n.wg.Done()
        for {
            conn, err := n.listener.Accept()
            if err != nil {
                select {
                case <-n.stopCh:
                    // we're shutting down
                    return
                default:
                    log.Printf("[%s] Accept error: %v", n.ID, err)
                    continue
                }
            }
            // Handle connection in a separate goroutine
            n.wg.Add(1)
            go func(c net.Conn) {
                defer n.wg.Done()
                n.handleConnection(c)
            }(conn)
        }
    }()
    return nil
}

// Stop stops listening and closes the node
func (n *Node) Stop() {
    close(n.stopCh)
    if n.listener != nil {
        _ = n.listener.Close()
    }
    n.wg.Wait()
}

// handleConnection handles incoming data from a Node or from a Client
func (n *Node) handleConnection(conn net.Conn) {
    defer conn.Close()

    buffer := make([]byte, 4096)
    for {
        // read data
        conn.SetDeadline(time.Now().Add(30 * time.Second))
        bytesRead, err := conn.Read(buffer)
        if err != nil {
            return
        }

        data := buffer[:bytesRead]

        // Parse the RelayCell from the raw data
        relayCell, err := ParseRelayCell(data)
        if err != nil {
            log.Printf("[%s] Error parsing relay cell: %v\n", n.ID, err)
            return
        }

        // If this is the exit node, we "finalize" the message (like sending to the destination)
        // For demonstration, we just log it
        if n.Role == ExitNode {
            log.Printf("[%s] EXIT node final output: %s\n", n.ID, relayCell.Payload)
            // In a real Tor network, you'd forward to the actual destination here.
        } else {
            // If not exit, forward to the next node
            // The relayCell.NextAddr is the IP:port of the next node
            if err := forwardToNextNode(relayCell); err != nil {
                log.Printf("[%s] Error forwarding to next node: %v\n", n.ID, err)
            }
        }
    }
}


// BootstrapNodes creates a list of Node instances from the config
// func BootstrapNodes(cfg *Config) []*Node {
//     nodes := make([]*Node, len(cfg.Ports))
//     for i, port := range cfg.Ports {
//         node := &Node{
//             ID:      fmt.Sprintf("Node%d", i+1),
//             Role:    RelayNode, // default to RELAY, roles will be reassigned upon path creation
//             Port:    port,
//             stopCh:  make(chan struct{}),
//         }
//         nodes[i] = node
//     }
//     return nodes
// }
